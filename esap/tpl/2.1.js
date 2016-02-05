var store = new Ext.data.JsonStore({
	autoLoad: true,
	fields: ['id', 'mType', 'mDesc', 'mUn', 'rem'],
	proxy: {
		type: 'ajax',
		url: '/esm/wmm',		
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
		{header: '备注', dataIndex: 'rem'},
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
	title: '物料表',
	id: 't2.1',	
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