var store = new Ext.data.JsonStore({
	autoLoad: true,
	fields: ['id', 'lcid', 'mType', 'mUn'],
	proxy: {
		type: 'ajax',
		url: '/esm/wmsl',		
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
		{header: '序号', dataIndex: 'id'},
		{header: '仓库', dataIndex: 'lcid'},
		{header: '物料分类', dataIndex: 'mType'},
		{header: '单位', dataIndex: 'mUn'},
		{
			text: '操作',
			xtype:'actioncolumn',
            width:50,
            items: [{
                icon: 'img/drop-add.gif',  // Use a URL in the icon config
                tooltip: 'add',
            },{
                icon: 'img/drop-yes.gif',
                tooltip: 'confirm',
            },{
                icon: 'img/drop-no.gif',
                tooltip: 'Delete', 
            }]
		}
	] 
}
var p = {
	xtype:'panel',
	title: '仓库信息',
	id: 't1.1',	
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
return p;