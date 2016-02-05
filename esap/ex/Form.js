Ext.define('ex.Form', {
    extend: 'Ext.form.Panel',
    alias: 'widget.xform',	

	winMod: true,
	form: null,
	cancelBtn: true,
	autoScroll: true,
	layout: 'anchor', 
	method: 'put',
	jsonSubmit: true,
	url: '',
	btnitems:[],
	cancelbtn:true,
	
    initComponent: function() {
		var me	= this;
		me.buttons=[];
		me.buttons = me.buttons.concat(me.btnitems);
		me.buttons.push({text: '保存', iconCls: 'icon-save', formBind: true, handler: me.onSubmitF});
		if(me.cancelbtn)me.buttons.push({text: '取消', handler: function(){this.up('window').close()}});
        this.callParent();
    },	
    onSubmitF: function() {
		var me = this.up('form');
		var fm = this.up('form').getForm();
		if (fm.isValid()) {			
			fm.submit({
				success: function(form, action) {
					Ext.Msg.alert('提示: 保存成功', action.result.msg);
					me.SubmitCallback();
				},
				failure: function(form, action) {
					Ext.Msg.alert('提示：保存失败', action.result.msg);
					me.SubmitCallback();
				}
			});
		}
	},
	
	SubmitCallback: function(){}
});