Ext.define('ex.Store', {
    extend: 'Ext.data.Store',
    alias: 'xstore',
    requires: [
        'Ext.data.proxy.Ajax',
        'Ext.data.reader.Json',
        'Ext.data.writer.Json'
    ],
    constructor: function(config) {
        config = Ext.apply({
            proxy: {
                type  : 'rest',
                reader: {type: 'json', root: 'data', totalProperty: 'total'},
                writer: 'json'
            }
        }, config);
		config.proxy.url = config.url
        this.callParent([config]);
    }
});