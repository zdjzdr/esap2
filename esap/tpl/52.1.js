var store = new Ext.data.JsonStore({
	autoLoad: true,
	fields: ['id', 'mType', 'mDesc', 'mUn', 'qty'],
	proxy: {
		type: 'ajax',
		url: '/esv/vZKC',
		reader: {
			root: 'data'
		}
	}
});
var g = {
	xtype: 'grid',
	store: store,
	columns: [
		new Ext.grid.RowNumberer(),
		// '',
		{header: '编码', dataIndex: 'id'},
		{header: '分类', dataIndex: 'mType'},
		{header: '描述', dataIndex: 'mDesc'},
		{header: '单位', dataIndex: 'mUn'},
		{header: '数量', dataIndex: 'qty'}
	],
	dockedItems: [{
        xtype: 'pagingtoolbar',
        store: store,
        dock: 'bottom',
        displayInfo: true
    }]
}
var p = {
	xtype:'panel',
	title: '总库存',
	id: 't52.1',	
	closable: true,
	tbar:[{
		text: "new",
		iconCls: 'up'
	},'->',{
		text: "login",
		iconCls: 'delete'
	},{
		text:"change",
		iconCls: 'add'
	}],
	items: [g]
}
return p;