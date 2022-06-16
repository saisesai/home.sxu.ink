export GLOBIGNORE=.gitignore:log/.gitignore

all: build

build: home.sxu.ink.exe
	$(MAKE) -C ./view all
	mv ./home.sxu.ink.exe ./build
	cp ./run.sh ./build
	cp -r ./public ./build/public

home.sxu.ink.exe:
	go build -ldflags="-s -w" -o $@ .

clean:
	rm -rf ./public/*
	rm -rf ./build/*
	mkdir ./build/data
	mkdir ./build/data/log
	rm -rf ./.idea/httpRequests/*
	rm -rf ./data/log/*