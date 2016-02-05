var cp = Ext.create('Ext.picker.Color', {
    value: '993300',  // 初始选择的颜色
    renderTo: Ext.getBody(),
    listeners: {
        select: function(picker, selColor) {
            alert(selColor);
        }
    }
});
var fp = {
	xtype: 'panel',
	id: 't61.1',
	title: '初始化',  
	closable: true,
	autoScroll: true,
	padding: 50,
	items: [cp]
}
return {
	c: fp
}