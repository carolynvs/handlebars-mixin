const helpers = require('/porter/mixins/handlebars/node_modules/handlebars-helpers')

module.exports = {
    register(handlebars) {
        helpers({ handlebars })
    }
}
