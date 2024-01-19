#!/usr/bin/env bash

function deploy_from_git_repo() {(
    # arguements
    BRANCH="$1"
    REPO="$2"
    CLONEDIR="$3"
    read -a BUILDDEPS <<< "$4"
    read -a DEPS <<< "$5"
    INSTALLFUNC="$6"

    # prep array of build dependencies to remove after successful build
    BUILDDEPS_TO_REMOVE=()

    # exit on error
    set -e

    # temporarily install build deps
    for pkg in "${BUILDDEPS[@]}"; do
        # check if package installed
        if ! dpkg -l "${PACKAGE}" > /dev/null 2>&1; then
            # package is not installed, install it and mark for uninstallation after build
            apt-get install --no-install-recommends -y "${pkg}"
            BUILDDEPS_TO_REMOVE+=("${pkg}")
        fi
    done

    # install deps
    for pkg in "${DEPS[@]}"; do
        apt-get install --no-install-recommends -y "${pkg}"
    done

    # clone repo
    git clone --depth 1 --branch "${BRANCH}" "${REPO}" "${CLONEDIR}"

    # run install tasks
    ldconfig
    $INSTALLFUNC

    # uninstall build deps that were installed for this build
    for pkg in "${BUILDDEPS_TO_REMOVE[@]}"; do
        apt-get remove -y "${pkg}"
    done
    apt-get autoremove -y

    # clean up
    rm -rf "${CLONEDIR}"
    
)}

function build_rtlsdr ()
{(
    set -e
    ldconfig
    mkdir -p "$CLONEDIR"/build
    pushd "$CLONEDIR"/build
    LD_LIBRARY_PATH="/usr/include" cmake \
        ../ \
        -DINSTALL_UDEV_RULES=ON \
        -DDETACH_KERNEL_DRIVER=ON \
        -DENABLE_ZEROCOPY=ON
    make
    make install
    ldconfig
    popd
)}

function build_mictronics_readsb ()
{(
    set -e
    ldconfig
    pushd "$CLONEDIR"
    make RTLSDR=yes HAVE_BIASTEE=yes
    mkdir -p /opt/Mictronics/readsb
    cp -v ./readsb.proto ./readsb ./readsbrrd ./viewadsb /opt/Mictronics/readsb
    ldconfig
    popd
)}

# ---------------------------------------------------------------------------

# install prerequisites
apt-get update -y
apt-get install --no-install-recommends -y isc-dhcp-client

# switch from NetworkManager to systemd-networkd
systemctl stop NetworkManager
systemctl disable NetworkManager
systemctl mask NetworkManager
systemctl unmask systemd-networkd
systemctl enable systemd-networkd

# remove existing netplan configuration
rm -v /etc/netplan/*.yml /etc/netplan/*.yaml

# deploy rtl-sdr
deploy_from_git_repo \
    "v2.0.1" \
    "https://gitea.osmocom.org/sdr/rtl-sdr.git" \
    "/src/rtl-sdr" \
    "build-essential cmake git libusb-1.0-0-dev pkgconf" \
    "" \
    build_rtlsdr

# deploy Mictronics readsb-protobuf
deploy_from_git_repo \
    "v4.0.2" \
    "https://github.com/Mictronics/readsb-protobuf.git" \
    "/stc/Mictronics-readsb-protobuf" \
    "build-essential git libncurses-dev libprotobuf-c-dev librrd-dev libusb-1.0-0-dev protobuf-c-compiler pkgconf" \
    "libncurses6 libprotobuf-c1" \
    build_mictronics_readsb
