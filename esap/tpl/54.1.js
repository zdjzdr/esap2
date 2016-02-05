var store = new Ext.data.JsonStore({
	autoLoad: true,
	fields: ['mid', 'mType', 'mDesc', 'mUn', 'gr','gi','mDate','lcid','lot','rem'],
	proxy: {
		type: 'ajax',
		url: '/esv/vJCMX',
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
		{header: '日期', dataIndex: 'mDate', xtype: 'datecolumn', format:'Y-m-d'},
		{header: '仓库', dataIndex: 'lcid'},
		{header: '物料', dataIndex: 'mid'},
		{header: '描述', dataIndex: 'mDesc'},		
		{header: '批号', dataIndex: 'lot'},
		{header: '入库数', dataIndex: 'gr'},
		{header: '出库数', dataIndex: 'gi'},
		{header: '单位', dataIndex: 'mUn'},
		{header: '分类', dataIndex: 'mType'},
		{header: '备注', dataIndex: 'rem'}
	],
	dockedItems: [{
        xtype: 'pagingtoolbar',
        store: store,
        dock: 'bottom',
        displayInfo: true
    }],
	listeners: {
		itemdblclick: function(me, rec){
			this.up('panel').down('textfield').setValue(rec.data.mid);
			fFn();
		}
	}
};
var fFn = function() {
	var val = Ext.getCmp(p.id + '-value').getValue();
	if (!val) {
		store.clearFilter();
		return;	
	}
	val = String(val).trim().split(" ");
	store.filterBy(function(rec) {
		var data = rec.data;
		for (var p in data) {
			var prop = String(data[p]);
			var len = val.length
			for (var i=0; i<len; i++) {
				//构造正则
				var matcher = val[i];
				var er = Ext.escapeRe;
				matcher = String(matcher);
				matcher = '^' + er(matcher);
				matcher = new RegExp(matcher);
				//遍历属性
				if (matcher.test(prop)) {
					return true;
				}
			}
		}
		return false;
	});
};
var p = {
	xtype:'panel',
	title: '进出明细',
	id: 't54.1',	
	closable: true,
	autoScroll: true,
	tbar:[{
		id: 't54.1-value',
		xtype: 'textfield',		 
		labelAlign: 'right',
		fieldLabel: "搜索",
		listeners: {
			specialkey: function(field, e) {
				// e.HOME, e.END, e.PAGE_UP, e.PAGE_DOWN,
				// e.TAB, e.ESC, arrow keys: e.LEFT, e.RIGHT, e.UP, e.DOWN
				switch(e.getKey()) {
					case e.ENTER:
						fFn();
						break;
					case e.ESC:
						// Ext.getCmp(p.id + '-value').setValue("");
						this.setValue("");
						store.clearFilter();
						break;
				}
			}
		}
	},{
		text: '查询',
		handler: fFn		
	},'-',{
		text: '清空',
		handler: function() {
			Ext.getCmp(p.id + '-value').setValue("");
			store.clearFilter();
		}
	},'->',{
		text: "?",
		handler: function() {
			Ext.Msg.alert('帮助', '友情提示：双击可以快速筛选物料编码。');
		}
	}],
	items: [g]
};
return p;