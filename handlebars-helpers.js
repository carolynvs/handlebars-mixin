const helpers = require('/opt/hbs/node_modules/handlebars-helpers')

module.exports = {
    register(handlebars) {
        helpers({ handlebars })
    }
}
