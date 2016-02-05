var store = Ext.create('Ext.data.Store', {
	fields:['id','name','email'],
	 proxy: {
		 type: 'ajax',
		 url: 'data/data.json',
		 reader: {
			 type: 'json',
			 root: 'data'
		 }
	 },
	 autoLoad: true
});
var p = {
	xtype: 'panel',
	id: 't14.1',
	title: '测试记录',  
	closable: true,
	autoScroll:true,
	frame: true,
	overflow:'auto',				
	items: [{
		xtype: 'grid',					
		title: '主表',
		store: store,
		columns: [
			{header:'id',dataIndex:'id',width:100},
			{header:'姓名',dataIndex:'name',width:75},	
			{header:'考勤记录',dataIndex:'email',width:150}	
		],
		listeners: {
			itemdblclick: function(v, rec) {
				alert(Ext.encode(rec.getId()));
			}
		},
		dockedItems:[{
			xtype:'pagingtoolbar',
			store: store,
			dock:'bottom',
			displayInfo:true,
			id:'dk1'
		}]
	}]
}
return {
	c:p
};