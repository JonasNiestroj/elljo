const fetch = require('node-fetch')
const tar = require('tar')
const zlib = require('zlib')
const mkdirp = require('mkdirp')
const fs = require('fs')
const exec = require('child_process').exec;
const path = require('path')

const binPath = path.join(__dirname, 'bin', 'jo');

function verifyAndPlaceBinary() {
    fs.renameSync("jo-compiler", binPath);
}

const install = async (callback) => {
    let gunzip = zlib.createGunzip();
    let untar = tar.x({path: "./bin"});

    // Pipe error to callback
    gunzip.on('error', callback);
    untar.on('error', callback);

    // Copy file to bin directory
    untar.on('end', verifyAndPlaceBinary);

    let req = await fetch("https://jo-compiler.s3.eu-central-1.amazonaws.com/jo_v0.0.0-SNAPSHOT-9e01fd1_darwin_amd64.tar.gz");
    req.body.pipe(gunzip).pipe(untar)
}

function uninstall(callback) {

}

let argv = process.argv;
if (argv && argv.length > 2) {
    let cmd = process.argv[2];
    if(cmd === "install") {
        install((err) => {
            if (err) {
                console.error(err);
                process.exit(1);
            } else {
                process.exit(0);
            }
        })
    } else if(cmd === "uninstall") {
        // TODO: Implement uninstall
        uninstall()
    }
}