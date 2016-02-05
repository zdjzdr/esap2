var store = new Ext.data.JsonStore({
	autoLoad: true,
    fields: ['id', 'lcid', 'mDate', 'ccid', 'cDate'],
	proxy: {
		type: 'ajax',
		url: '/esm/wmgi',		
		reader: {
			root: 'data',
		}
	}
});
var g = {
	xtype: 'grid',
	store: store,
	columns: [
		new Ext.grid.RowNumberer(),
		{header: 'mDate', dataIndex: 'mDate'},
		{header: 'lcid', dataIndex: 'lcid'},
		{header: 'ccid', dataIndex: 'ccid'},
		{header: 'id', dataIndex: 'id'},
		{header: 'cdate', dataIndex: 'cDate'},
	] 
}
var p = {
	xtype:'panel',
	title: '出库单',
	id: 't522.1',	
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
return {
	c: p
};