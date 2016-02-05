var fp = {
	xtype:'form',
	title: '员工信息',
	id: 't121.1',
	closable: true,
	frame: true,
	items: [{
		xtype: 'datefield',
		name: 'startDate',
		format: 'Y年m月d日',
		disableDates: ['10日'],
		fieldLabel:'起始日'
	}]
};
return {
	c:fp
}