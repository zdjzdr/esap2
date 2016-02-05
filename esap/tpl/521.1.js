var store = new Ext.data.JsonStore({
	autoLoad: true,
    fields: ['id', 'lcid', 'mDate', 'vid', 'cDate'],
	proxy: {
		type: 'ajax',
		url: '/esm/wmgr',		
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
		{header: 'vid', dataIndex: 'vid'},
		{header: 'id', dataIndex: 'id'},
		{header: 'cdate', dataIndex: 'cDate'},
	] 
}
var p = {
	xtype:'panel',
	title: '入库单',
	id: 't521.1',	
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