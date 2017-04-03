# Holds short targets for the go tool commands.

# Watch for changes, build, and push.
watch:
	watchman-make -p 'src/**/*.go' -t push

run:
	go get juju.net.nz/pipoint/pipoint
	./bin/pipoint

check:
	go test juju.net.nz/pipoint

push:
	GOARCH=arm GOARM=7 go get juju.net.nz/pipoint/pipoint
	rsync -zt bin/linux_arm/pipoint pi-ed7:~
