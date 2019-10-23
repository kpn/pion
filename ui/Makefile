# Build app that uses local APIs, i.e. in memory services
nodes_modules: 
	npm install

.PHONY: build
build: nodes_modules
	ng build

# Build app that optimized for production
.PHONY: build-prod
build-prod: nodes_modules
	ng build -c production

# Build app 'remote' profile, i.e. no optimization, calling remote API endpoints
.PHONY: build-remote
build-remote: nodes_modules
	ng build -c remote

.PHONY: clean
clean:
	rm -rf dist/
