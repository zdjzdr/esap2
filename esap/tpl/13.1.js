var fp = {
	xtype: 'form',
	id: 't13.1',
	title: '考勤记录',  
	closable: true,
	labelWidth: 90,
	items: [
		{
			xtype: 'fieldset',
			title: '折叠项',
			autoHeight: true,
			defaultType: 'checkbox',
			collapsible: true,
			items: [
				{fieldLabel:'项目1'},
				{fieldLabel:'项目2'},
				{fieldLabel:'项目3'},
				{fieldLabel:'项目4'}
			]
		}
	]
}
return {
	c: fp
}