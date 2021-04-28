const request = require('request'),
      tar = require('tar'),
      mkdirp = require('mkdirp');

function install() {
    mkdirp.sync("./bin");
    let ungz = zlib.createGunzip();
    let untar = tar.Extract({ path: "./bin" });

    ungz.on('error', callback);
    untar.on('error', callback);

    // First we will Un-GZip, then we will untar. So once untar is completed,
    // binary is downloaded into `binPath`. Verify the binary and call it good
    untar.on('end', verifyAndPlaceBinary.bind(null, "jo", "./bin", callback));

    console.log("Downloading from URL: https://jo-compiler.s3.eu-central-1.amazonaws.com/jo_v0.0.0-SNAPSHOT-9e01fd1_darwin_amd64.tar.gz");
    let req = request({ uri: "https://jo-compiler.s3.eu-central-1.amazonaws.com/jo_v0.0.0-SNAPSHOT-9e01fd1_darwin_amd64.tar.gz" });
    req.on('error', callback.bind(null, "Error downloading from URL: https://jo-compiler.s3.eu-central-1.amazonaws.com/jo_v0.0.0-SNAPSHOT-9e01fd1_darwin_amd64.tar.gz"));
    req.on('response', function (res) {
        if (res.statusCode !== 200) return callback("Error downloading binary. HTTP Status Code: " + res.statusCode);
    });
}

function uninstall(callback) {}

let argv = process.argv;
if (argv && argv.length > 2) {
    let cmd = process.argv[2];
    if (cmd === "install") {
        console.log("tet");
        install(function (err) {
            if (err) {
                console.error(err);
                process.exit(1);
            } else {
                console.log("INSTALLED");
                process.exit(0);
            }
        });
    } else if (cmd === "uninstall") {
        uninstall();
    }
}