var p = {
	xtype:'form',
	title: '员工培训',
	id: 't122.1',
	closable: true,
	frame: true,
	layout: 'absolute',
	defaults: {
		xtype: 'label'
	},
	items: [
		{
			x: 100,
			y: 50,
			html: '姓名<font color="red">>></font>',
		},{
			x: 150,
			y: 45,
			width: 120,
			id: 'userName',
			hiddenName: 'userName',
			xtype: 'textfield'
		},{
			x: 285,
			y: 50,
			html: '<font color="red">*</font>',
		},{
			x: 100,
			y: 80,
			width:60,
			html: '年龄<font color="red">>></font>',
		},{
			x: 150,
			y: 75,
			width: 120,
			id: 'age',
			name: 'age',
			hiddenName: 'age',
			xtype: 'textfield'
		},{
			x: 285,
			y: 80,
			html: '<font color="red">*</font>'
		}
	]
};
return {
	c:p
}