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
const { exec } = require('child_process')
/**
 * @type {Cypress.PluginConfig}
 */
// eslint-disable-next-line no-unused-vars
module.exports = (on, config) => {
    on('task', {
        async execJo(options) {
            await exec('../../../main ../fixtures/' + config.input + ' ../fixtures/' + config.index + ' ../' + config.output)
        }
    })
  // `on` is used to hook into various events Cypress emits
  // `config` is the resolved Cypress config
}
