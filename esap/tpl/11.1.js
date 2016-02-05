var store = new Ext.data.JsonStore({
	autoLoad: true,
    fields: ['id', 'name', 'lng', 'lat'],
	proxy: {
		type: 'ajax',
		url: '/esm/wmm',		
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
		// '',
		{header: 'id', dataIndex: 'id'},
		{header: 'name', dataIndex: 'name'},
		{header: 'lng', dataIndex: 'lng'},
		{header: 'lat', dataIndex: 'lat'},
	] 
}
var p = {
	xtype:'panel',
	title: '待办事宜',
	id: 't11.1',	
	closable: true,
	// frame: true,
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