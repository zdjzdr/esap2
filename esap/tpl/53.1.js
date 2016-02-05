var store = new Ext.data.JsonStore({
	autoLoad: true,
	fields: ['mid', 'mType', 'mDesc', 'mUn', 'qty', 'lot', 'lcid'],
	proxy: {
		type: 'ajax',
		url: '/esv/vPCKC',
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
		{header: '仓库', dataIndex: 'lcid'},
		{header: '编码', dataIndex: 'mid'},
		{header: '分类', dataIndex: 'mType'},
		{header: '描述', dataIndex: 'mDesc'},
		{header: '批号', dataIndex: 'lot'},
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
	title: '批次库存',
	id: 't53.1',	
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