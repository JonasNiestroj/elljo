const fetch = require('node-fetch');
const tar = require('tar');
const zlib = require('zlib');
const mkdirp = require('mkdirp');
const fs = require('fs');
const exec = require('child_process').exec;
const path = require('path');

function getInstallationPath(callback) {

    // `npm bin` will output the path where binary files should be installed
    exec("npm bin", function (err, stdout, stderr) {

        let dir = null;
        if (err || stderr || !stdout || stdout.length === 0) {

            // We couldn't infer path from `npm bin`. Let's try to get it from
            // Environment variables set by NPM when it runs.
            // npm_config_prefix points to NPM's installation directory where `bin` folder is available
            // Ex: /Users/foo/.nvm/versions/node/v4.3.0
            let env = process.env;
            if (env && env.npm_config_prefix) {
                dir = path.join(env.npm_config_prefix, "bin");
            }
        } else {
            dir = stdout.trim();
        }

        mkdirp.sync(dir);

        callback(null, dir);
    });
}

function verifyAndPlaceBinary(binName, binPath, callback) {
    fs.renameSync(path.join(binPath, binName), path.join(installationPath, binName));
    getInstallationPath(function (err, installationPath) {
        if (err) return callback("Error getting binary installation path from `npm bin`");

        // Move the binary file
        fs.renameSync(path.join(binPath, binName), path.join(installationPath, binName));

        callback(null);
    });
}

const install = async callback => {
    mkdirp.sync("./bin");

    let gunzip = zlib.createGunzip();
    let untar = tar.x({ path: "./bin" });

    // Pipe error to callback
    gunzip.on('error', callback);
    untar.on('error', callback);

    // Copy file to bin directory
    //untar.on('end', verifyAndPlaceBinary.bind(null, "jo", "./bin", callback));

    let req = await fetch("https://jo-compiler.s3.eu-central-1.amazonaws.com/jo_v0.0.0-SNAPSHOT-9e01fd1_darwin_amd64.tar.gz");
    req.body.pipe(gunzip).pipe(untar);
};

function uninstall(callback) {}

let argv = process.argv;
if (argv && argv.length > 2) {
    let cmd = process.argv[2];
    if (cmd === "install") {
        install(err => {
            if (err) {
                console.error(err);
                process.exit(1);
            } else {
                process.exit(0);
            }
        });
    } else if (cmd === "uninstall") {
        // TODO: Implement uninstall
        uninstall();
    }
}