Ext.define('ex.Grid2', {
    extend: 'Ext.grid.Panel',
    alias: 'widget.xgrid2',

    requires: [
        'Ext.grid.plugin.CellEditing',
        'Ext.toolbar.TextItem'
    ],
	
	autoScroll: true, 
	maxHeight: 320,
	store: null,
	dock: 'top',
	selType: 'cellmodel',
	selecRec: null,
	tbar: [],
	plugins: [],
	columns: [],
	
    initComponent: function() {
		var me = this;
		me.store.load();
        me.editing = Ext.create('Ext.grid.plugin.CellEditing', {clicksToEdit: 1});
        me.plugins.push(me.editing);
		
		me.tbar.push({
			itemId:'tbar_add', text: '添加', iconCls: 'icon-add',
			handler: function() {
					var rec = {};
					me.store.insert(0, rec);
					me.editing.startEditByPosition({row: 0, column: 0});
			}
		});
		
		me.columns.push({
			xtype:'actioncolumn', text: 'Go', width:50,
			items: [{
				iconCls: 'icon-add', tooltip: '插入(行后)',
				handler: function(grid, rowIndex, colIndex) {
					var rec = {};
					me.getStore().insert(rowIndex + 1, rec);
					me.editing.startEditByPosition({row: rowIndex + 1, column: 0});
				}
			},{
				iconCls: 'icon-del', tooltip: '删除',
				handler: function(grid, rowIndex, colIndex) {
					me.getStore().removeAt(rowIndex);
				}
			}]
		});
        me.callParent();
    }
});