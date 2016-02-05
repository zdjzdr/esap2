/**
 * 模板变量
 * 字段、数据表、模板名称、模板id, grid表头
 * by woylin , 2015/11/27
 */
// var test = new Ext.create('js.aa',{});
var c1 = ['mDate', 'lcid', 'vid', 'id', 'cDate', 'ExcelServerRCID'],
	c2 = ['mid', 'lot', 'qty', 'rem', 'id'],
	m1 = 'wmgr',
	m2 = '' + m1 + '_d',
	tt1 = '入库单',
	tt2 = '入库明细',
	id1	= 't3.1';
var g1h1 = {header: '日期', dataIndex: 'mDate', xtype: 'datecolumn', format:'Y-m-d'},
	g1h2 = {header: '仓库', dataIndex: 'lcid'},
	g1h3 = {header: '供方', dataIndex: 'vid'},
	g1h4 = {header: '单号', dataIndex: 'id'},
	g1h5 = {header: '结算日', dataIndex: 'cDate'};
var g2h1 = {header: '编号', dataIndex: 'mid', editor: 'textfield'},
	g2h2 = {header: '批号', dataIndex: 'lot', editor: 'textfield'},
	g2h3 = {header: '数量', dataIndex: 'qty', editor: 'numberfield'},
	g2h4 = {header: '备注', dataIndex: 'rem', editor: 'textfield'};
//数据仓库定义store for grid1,grid2
var reader = {root: 'data', totalProperty: 'total'};
var s1 = new Ext.data.JsonStore({
	autoLoad: true, fields: c1,
	proxy: {type: 'rest', reader: reader, url: '/esm?m=' + m1}
});
var s2 = new Ext.data.JsonStore({
	fields: c2,
	proxy: {type: 'rest', reader: reader, url: '/esd?m=' + m2}
});
//主表单双击事件listners for grid1
var gClk = function(me, rec, index) {
	var id = rec.data.ExcelServerRCID;
	g2.store.proxy.url = '' + '/esd?m=' + m2 + '&rcid=' + id;
	g2.store.loadPage(1);
};
var gDblclk = function(me, rec, item, index) {	
	var submitF = function(){
		fm2.store.sync();
		var form = fm.getForm();
		// fm.getForm().setValues({data: storeToJson(fm2.getStore())});
		if (form.isValid()) {
			form.submit({
				success: function(form, action) {
					Ext.Msg.alert('保存成功', action.result.msg);
					s1.load(),s2.load();
					win.close();
				},
				failure: function(form, action) {
					Ext.Msg.alert('保存失败', action.result.msg);
					win.close();
				}
			});
		}
	};
	var rcid = rec.data.ExcelServerRCID;
	var fm2s = new Ext.data.JsonStore({
		fields: c2, autoLoad: true,
		proxy: {type: 'rest', reader: reader, url: '/esd?m=' + m2 + '&rcid=' + rcid}
	});	
	var fm2e = Ext.create('Ext.grid.plugin.CellEditing', {clicksToEdit: 1});
	var fm2 = Ext.create('Ext.grid.Panel',{
		xtype: 'grid', store: fm2s, autoScroll: true, selType: 'cellmodel', maxHeight: 200,
		plugins: [fm2e],
		columns: [
			// new Ext.grid.RowNumberer(),
			g2h1,g2h2,g2h3,g2h4, 		//表头2
			{
				xtype:'actioncolumn', text: 'Go', width:50,
				items: [{
					iconCls: 'icon-add', tooltip: '插入(行后)',
					handler: function(grid, rowIndex, colIndex) {
						var rec = {};
						fm2s.insert(rowIndex + 1, rec);
						fm2e.startEditByPosition({row: rowIndex + 1, column: 0});
					}
				},{
					iconCls: 'icon-del', tooltip: '删除',
					handler: function(grid, rowIndex, colIndex) {
						fm2s.removeAt(rowIndex);
					}
				}]
			}
		],
		tbar: [{
                text: '添加', iconCls: 'icon-add',
                handler: function() {
						var rec = {};
						fm2s.insert(0, rec);
						fm2e.startEditByPosition({row: 0, column: 0});
				}
        }]
	});
	var fm = Ext.create('Ext.form.Panel', {
		autoScroll: true,layout: 'anchor',
		method: 'put', jsonSubmit: true,
		url: 'esm?m=' + m1 + '&rcid=' + rec.data.ExcelServerRCID,
		defaults:{xtype: 'textfield', labelAlign: 'right'},
		buttons: [
			'->',{
				text: '保存', iconCls: 'icon-save', handler: submitF
			},{
				text: '取消', handler: function() {win.close();}
			}
		],
		items: [{
			xtype: 'fieldset', title:'主表', margin: '0 5 5 5', padding: 5, defaultType: 'textfield',
			items: [{
					xtype: 'datefield',
					name: 'mDate', fieldLabel:'日期', anchor: '100%', format: 'Y-m-d',
				},{	
					name: 'lcid', fieldLabel: '仓库', anchor: '100%'
				},{	
					name: 'vid', fieldLabel: '供方', anchor: '100%'
				},{	
					name: 'id', fieldLabel: '单号', anchor: '100%', xtype: 'displayfield'
				},{	
					name: 'data', fieldLabel: '数据', xtype: 'hidden'
				},{	
					name: 'cDate', fieldLabel: '结算日', xtype: 'displayfield'
				},
				fm2
			]
		}],
	});
	fm.getForm().loadRecord(rec);
	
	var win = Ext.create('Ext.window.Window', {
		autoScroll: true, title: tt1, width: 500, items: [fm], tools: [{type: 'pin'}]
	}).show();
};
//主表grid1
var g = {
	xtype: 'grid', store: s1, selType: 'checkboxmodel',
	columns: [
		new Ext.grid.RowNumberer(),
		g1h1,g1h2,g1h3,g1h4,g1h5,	//表头1
		{
			xtype:'actioncolumn', width:50,
            items: [{
                iconCls: 'icon-yes', tooltip: '确认收货',
                handler: function(grid, rowIndex, colIndex) {
                    var rec = grid.getStore().getAt(rowIndex);
                    alert("确认收货 " + rec.get('id'));
                }
            },{
                iconCls: 'icon-del', tooltip: '删除',
                handler: function(grid, rowIndex, colIndex) {
                    var rec = grid.getStore().getAt(rowIndex);
                    alert("删除 " + rec.get('id'));
                }
            }]
		}
	],
	dockedItems: [{xtype: 'pagingtoolbar', store: s1, dock: 'top', displayInfo: true}],
	listeners: {
		beforeselect: gClk,
		itemdblclick: gDblclk
	}
}
//明细表grid2
var g2 = {
	xtype: 'grid', autoScroll: true,store: s2, maxHeight: 200,
	columns: [
		new Ext.grid.RowNumberer(),
		g2h1,g2h2,g2h3,g2h4			//表头2
	],
	dockedItems: [{xtype: 'pagingtoolbar', store: s2, dock: 'bottom', displayInfo: true}]
};
//闭包panel容器
var p = {
	xtype:'panel', title: tt1, id: id1,	glyph: 99, autoScroll: true, closable: true,
	tbar:[{
		text: "新建", iconCls: 'icon-add'
	}],
	items: [g, g2]
}
return p;