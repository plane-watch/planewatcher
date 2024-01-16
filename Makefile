help:
	$(info )
	$(info Valid make targets are:)
	$(info   - uefi-x86-ovf : Outputs a zip file containing .vmx & .vmdk for uefi-x86 platform)
	$(info )

uefi-x86-ovf:
	# Clone armbian/build repo if needed
	if [ ! -d "./armbian-build" ]; then git clone --depth=1 --branch=v23.11 https://github.com/armbian/build ./armbian-build; fi
	# Copy armbian customisations
	cp -Rv ./armbian/* ./armbian-build/
	# Run armbian build
	cd ./armbian-build && ./compile.sh build BOARD=uefi-x86 ENABLE_EXTENSIONS=image-output-ovf planewatcher

clean:
	# Remove armbian/build
	rm -r ./armbian-build
