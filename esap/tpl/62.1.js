var fp = {
	xtype: 'form',
	id: 't62.1',
	title: '全局设置',  
	closable: true,
	autoScroll: true,
	labelWidth: 90,
	padding: 10,
	defaults: {xtype: 'textfield', labelAlign: 'right'},
	items: [
		{
			xtype: 'fieldset',
			title: '同时应用到',
			autoHeight: true,
			// defaultType: 'checkbox',
			// defaultAlign: 'right',
			defaults: {xtype: 'checkbox', labelAlign: 'right'},
			collapsible: true,
			items: [
				{fieldLabel:'人力'},
				{fieldLabel:'财务'},
				{fieldLabel:'销售'},
				{fieldLabel:'生产'},
				{fieldLabel:'仓库'}
			]			
		},
		{fieldLabel: '单价'},
		{fieldLabel: '存货单重'},
		{fieldLabel: '成本'},
		{fieldLabel: 'BOM单位用量'},
		{fieldLabel: '单位转换率'},
		{fieldLabel: '税率'},
		{fieldLabel: '百分比'},
		{fieldLabel: '小时'},
		{fieldLabel: '秒'},
		{fieldLabel: '箱数'},
		{xtype: 'label', html: '<font color="red">一经设置，不得更改</font>'},
		{fieldLabel: '系统超时'}
	]
}
return {
	c: fp
}