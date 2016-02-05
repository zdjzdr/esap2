var store = new Ext.data.JsonStore({
	autoLoad: true,
    fields: ['id', 'name', 'mType', 'mUn'],
	proxy: {
		type: 'ajax',
		url: '/ajax?m=wmm',		
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
		{header: 'mtype', dataIndex: 'mType'},
		{header: 'mUn', dataIndex: 'mUn'},
	] 
}
var p = {
	xtype:'panel',
	title: '物料表',
	id: 't51.1',	
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