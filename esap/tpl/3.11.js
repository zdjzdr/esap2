/**
 * 模板变量
 * 字段、数据表、模板名称、模板id, grid表头
 * by woylin , 2015/11/27
 */
Ext.require(['js.wm.jForm', 'js.wm.jWin']);
var c1 = ['mDate', 'lcid', 'vid', 'id', 'cDate', 'ExcelServerRCID'],
	c2 = ['mid', 'lot', 'qty', 'rem', 'id'],
	m1 = 'wmgr',
	m2 = '' + m1 + '_d',
	tt1 = '入库单',
	tt2 = '入库明细',
	id1	= 't3.1';
var g1h1 = {header: '日期', dataIndex: 'mDate', editor: 'textfield', xtype: 'datecolumn', format:'Y-m-d'},
	g1h2 = {header: '仓库', dataIndex: 'lcid', editor: 'textfield'},
	g1h3 = {header: '供方', dataIndex: 'vid', editor: 'textfield'},
	g1h4 = {header: '单号', dataIndex: 'id', renderer: renderId},
	g1h5 = {header: '结算日', dataIndex: 'cDate', renderer: renderDate};
var g2h1 = {header: '编号', dataIndex: 'mid', editor: 'textfield'},
	g2h2 = {header: '批号', dataIndex: 'lot', editor: 'textfield'},
	g2h3 = {header: '数量', dataIndex: 'qty', editor: 'numberfield'},
	g2h4 = {header: '备注', dataIndex: 'rem', editor: 'textfield'};
/**
 * 模板主体
 * 主从表模板：主表-明细 支持增删改
 * by woylin , 2015/11/27
 */
//数据仓库定义store for grid1,grid2
var s1 = new Ext.data.JsonStore({
	autoLoad: true, fields: c1,
	proxy: {type: 'rest', reader: reader, url: '/esm/' + m1 + '?'}
});
var s1sync = function() {
	s1.sync();
	s1.load();
	//this.up('panel').down('#saveBtn').disable();
};
var s2 = new Ext.data.JsonStore({
	fields: c2,
	proxy: {type: 'rest', reader: reader, url: '/esd/' + m2 + '?'}
});
//*************修改窗口HEAD***************
var fm2s = new Ext.data.JsonStore({
	fields: c2,
	proxy: {type: 'rest', reader: reader, url: '/esd/' + m2 + '?'}
});	
var fm2e = Ext.create('Ext.grid.plugin.CellEditing', {clicksToEdit: 1});
var fm2 = Ext.create('Ext.grid.Panel',{
	xtype: 'grid', store: fm2s, autoScroll: true, selType: 'cellmodel', maxHeight: 200,
	plugins: [fm2e],
	columns: [
		new Ext.grid.RowNumberer(),
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
var submitF = function(){
	fm2.store.sync();
	var form = fm.getForm();
	if (form.isValid()) {
		form.submit({
			success: function(form, action) {
				Ext.Msg.alert('提示', '保存成功');
				s1.load(), s2.load();
				win.close();
			},
			failure: function(form, action) {
				Ext.Msg.alert('提示：保存失败', action.result.msg);
				win.close();
			}
		});
	}
};
var fm = Ext.create('Ext.form.Panel', {
	autoScroll: true,layout: 'anchor',
	method: 'put', jsonSubmit: true,
	url: 'esm/' + m1 + '?id=' + rec.data.id,
	items: [{
		xtype: 'fieldset', title:'主表', margin: '0 5 5 5', padding: 5,
		defaults: {xtype: 'textfield', labelAlign : 'right'},
		items: [
			{name: 'mDate', fieldLabel:'日期', format: 'Y-m-d', anchor: '100%'},
			{name: 'lcid', fieldLabel: '仓库', anchor: '100%'},
			{name: 'vid', fieldLabel: '供方', anchor: '100%'},
			{xtype: 'container', layout: 'hbox', defaults: {xtype: 'textfield', labelAlign : 'right'},
				items:[
					{name: 'id', fieldLabel: '单号', xtype: 'displayfield', renderer: renderId},
					{name: 'cDate', fieldLabel: '结算日', xtype: 'displayfield'}
				]
			},
			fm2
		]
	}],
	buttons: ['->',
		{text: '保存', iconCls: 'icon-save', handler: submitF},
		{text: '取消', handler: function() {win.close();}}
	],
});

//*************修改窗口END***************
//主表单击事件for grid1
var gClk = function(me, rec, index) {
	var id = rec.data.id;
	g2.store.proxy.url = '' + '/esd/' + m2 + '?id=' + id;
	fm2s.proxy.url = '' + '/esd/' + m2 + '?id=' + id;
	fm.getForm().loadRecord(rec);
	g2.store.loadPage(1);	
};
//主表grid1
// var g1e = Ext.create('Ext.grid.plugin.CellEditing', {clicksToEdit: 1});
var g1e = Ext.create('Ext.grid.plugin.RowEditing', {
	listeners: {
		cancelEdit: function(rowEditing, context) {
			if (context.record.phantom) {
				s1.remove(context.record);
			}
		}
	}
});
var g1 = Ext.create('Ext.grid.Panel', {
	xtype: 'grid', store: s1, selType: 'checkboxmodel',
	plugins: [g1e],
	columns: [
		new Ext.grid.RowNumberer(),
		g1h1,g1h2,g1h3,g1h4,g1h5,		//表头1
		{
			xtype:'actioncolumn', width:50,
            items: [{
                iconCls: 'icon-yes', tooltip: '过账',
                handler: function(grid, rowIndex, colIndex) {
                    var rec = grid.getStore().getAt(rowIndex);
                    alert("过账 " + rec.get('id'));
					rec.set('cDate', Ext.Date.format(new Date(), 'Y-m-d H:i:s'));
					s1sync();
                }
            },{
                iconCls: 'icon-del', tooltip: '删除',
                handler: function(grid, rowIndex, colIndex) {
                    Ext.Msg.alert("提示", '确认删除?', 
						function(btn){
							if (btn=='ok') {
								s1.removeAt(rowIndex);
								s1sync();
							}
						}
					);				
                }
            }]
		}
	],
	dockedItems: [{xtype: 'pagingtoolbar', store: s1, dock: 'top', displayInfo: true}],
	listeners: {
		beforeselect: gClk,
		// itemdblclick: gDblclk,
		beforeitemcontextmenu: rightClickFn
		
	}
});
//为右键菜单添加事件监听器
var rightClick = new Ext.menu.Menu({ 
    id:'rightClickCont', 
    items: [ 
        {   
			id: 'rMenu1', 
            iconCls:'ico_add',
            text: '添加数据', 
        },{ 
			id: 'rMenu2', 
            text: '修改数据',iconCls:'ico_search',
            handler:function() {
               var row = g1.getSelectionModel().getSelection();
               var id=row[0].get('id');
               alert(id);
            },
        },{
           text:'删除数据',iconCls:'ico_del'
        },{text:'导出数据'}
    ] 
}); 
//右键菜单代码关键部分 
function rightClickFn(me, rec, item, index, e){ 
    e.preventDefault(); 
    rightClick.showAt(e.getXY()); 
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
		text: "新建", iconCls: 'icon-add', 
		handler: function() {
			var rec = {mDate: Ext.Date.format(new Date(), 'Y-m-d H:i:s')};
			s1.insert(0, rec);
			// g1e.startEditByPosition({row: 0, column: 1});
			g1e.startEdit(0, 0);
			
		}
	},{
		text: '修改', iconCls: 'icon-yes',
		handler: function () {
			var selecRec = g1.getView().getSelectionModel().getSelection()[0];
			if (!selecRec) {
				Ext.Msg.alert("提示", '请先选择后再试。');
				return
			};
			Ext.create('Ext.window.Window', {
				autoScroll: true, title: tt1, width: 500, items: []
			}).show();
		}
	},{
		text: '删除', iconCls: 'icon-del',
		handler: function() {
			Ext.Msg.alert("提示", '确认删除?', 
				function(btn){
					if (btn=='ok') {
						var selecRec = g1.getView().getSelectionModel().getSelection()[0];
						if (selecRec) {s1.remove(selecRec);}
						s1sync();
					}
				}
			);	
		}
	},'-',{
		text: '保存', itemId: 'saveBtn',
		iconCls: 'icon-save',	handler: s1sync
	}],
	items: [g1, g2]
}
return p;