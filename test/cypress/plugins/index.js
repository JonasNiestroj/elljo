/// <reference types="cypress" />
// ***********************************************************
// This example plugins/index.js can be used to load plugins
//
// You can change the location of this file or turn off loading
// the plugins file with the 'pluginsFile' configuration option.
//
// You can read more here:
// https://on.cypress.io/plugins-guide
// ***********************************************************

// This function is called when a project is opened or re-opened (e.g. due to
// the project's config changing)
const fs = require('fs')
const path = require('path')
/**
 * @type {Cypress.PluginConfig}
 */
// eslint-disable-next-line no-unused-vars
module.exports = (on, config) => {
    on('task', {
        async spawnEllJo({ filePath }) {
            return new Promise((resolve, reject) => {
                var source = fs.readFileSync(filePath, 'utf-8').toString()
                var spawn = require('child_process').spawn
                const child = spawn(path.join(__dirname, '../../', 'main'), ['--service']);
                var command = `compile ${source.replace(/\r?\n|\r/g, "\\n")}`
                var buffer = Buffer.from(command, 'utf8')
                let outputJson = ""
                child.stdin.write(buffer);
                child.stdin.end();
                child.stdout.on('data', function (data) {
                    outputJson += data.toString()
                });
                child.on('close', () => {
                    let output = JSON.parse(outputJson)
                    let outputJs = output.output.replace("export default component", "")
                    const testHtml = fs.readFileSync('./test.html').toString()
                    resolve(testHtml.replace("{{SCRIPT}}", outputJs))
                })
            })
        }
    })
  // `on` is used to hook into various events Cypress emits
  // `config` is the resolved Cypress config
}
