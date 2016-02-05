Ext.define('ex.Tree', {
	extend: 'Ext.tree.Panel', 
	xtype: 'xtree',
	rootVisible: false,
	itemclk: function(v, rec) {
		if (rec.data.leaf) (iTpl = Es.chkTab(rec)) ? iTpl.show() : Es.addTab(rec);
	},
	initComponent: function() {
		if (!this.store) this.store = Es.loadTreeStore('data/' + this.itemId.substr(2) + 'Menu.json');
		this.listeners = {
			itemclick: this.itemclk,
			beforeitemcontextmenu: this.onRightClickFn
		};
		
		this.callParent();
	},
	onRightClick: new Ext.menu.Menu({
		// id:'treerightClickCont', 
		items: [{
				text: '新建', 
				iconCls:'icon_add',
			},{
				text: '刷新',
				iconCls:'icon_fresh',
			},{
				text:'导出excel数据清单',
				iconCls:'icon_save'
			},{
				text:'导出txt数据清单',
				iconCls:'icon_save'
		}]
	}),
	onRightClickFn: function(me, rec, item, index, e){
		e.preventDefault(); 
		this.onRightClick.showAt(e.getXY()); 
	}
});