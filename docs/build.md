# Build

## Requirements
Install the following packages :
- npm
	```bash
	brew install npm
	```
- angular-cli
	```bash
	brew install angular-cli
	```
## Build instructions
A Make file is provided at root directory. Some profiles are as follows:

- Clean, build and publish images:
	```bash
	make clean image
	```
- Clean and build binary files:
	```bash
	make clean build
	```
- Run unit-tests
	```bash
	make test
	```
## Development mode
To enable UI component in development mode, you can set the environment variable `BUILD_ENV=dev`, e.g.:
```bash
export BUILD_ENV=dev
make clean image
```