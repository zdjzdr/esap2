Ext.define('ex.Grid', {
    extend: 'Ext.grid.Panel',
    alias: 'widget.xgrid',

    requires: [
        'Ext.grid.plugin.CellEditing',
        'Ext.form.field.Text',
        'Ext.toolbar.TextItem'
    ],
	
	// editing:{},
	store: null,
	dock: 'top',
	selType: 'checkboxmodel',
	selecRec: null,
	plugins: [], 
	dockedItems: [],
	
    initComponent: function() {
		var me = this;
        me.editing = Ext.create('Ext.grid.plugin.RowEditing', {
			listeners: {
				cancelEdit: function(rowEditing, context) {
					if (context.record.phantom) {
						me.store.remove(context.record);
					}
				}
			}
		});
		me.plugins.push(me.editing);
		
		this.dockedItems.push({
			xtype: 'pagingtoolbar', 
			store: this.store, 
			dock: this.dock, 
			displayInfo: true,
			itemId:'pagebar_1'
		});	
		this.listeners = {
			beforeselect: me.onClick,
			// itemdblclick: gDblclk,
			beforeitemcontextmenu: me.onRightClickFn			
		};

        this.callParent();
    },
    onClick: function(me, rec, index) {},
	onRightClick: new Ext.menu.Menu({
		// id:'gridrightClickCont', 
		items: [{
				id: 'rMenu1', 
				iconCls:'ico_add',
				text: '添加数据', 
			},{
				id: 'rMenu2', 
				text: '修改数据',iconCls:'ico_search',
				handler:function() {
				   var row = g1.getSelectionModel().getSelection();
				   var id=row[0].get('id');
				   alert(id);
				},
			},{
			   text:'删除数据',iconCls:'ico_del'
			},{text:'导出数据'}
		] 
	}),
	onRightClickFn: function(me, rec, item, index, e){
		e.preventDefault(); 
		this.onRightClick.showAt(e.getXY()); 
	},
	startEdit: function(rec, col) {
		this.editing.startEdit(rec, col);
	},
	getSelectRec: function() {
		this.selecRec = this.getSelectionModel().getSelection()[0];
		if (!this.selecRec) {
			Ext.Msg.alert("提示", '请先选择后再试。');
			return null
		};
		return this.selecRec;		
	}
});