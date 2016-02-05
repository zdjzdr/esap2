/**
 * 模板变量
 * 字段、数据表、模板名称、模板id, grid表头
 * by woylin , 2015/11/27
 */
var c1 = ['mDate', 'lcid', 'ccid', 'id', 'cDate', 'ExcelServerRCID'],
	c2 = ['mid', 'lot', 'qty', 'rem', 'id'],
	m1 = 'wmgi',
	m2 = '' + m1 + '_d',
	h1 = '出库单',
	t1 = 't4.1';
var g1h1 = {header: '日期', dataIndex: 'mDate', editor: 'textfield', xtype: 'datecolumn', format:'Y-m-d'},
	g1h2 = {header: '仓库', dataIndex: 'lcid', editor: 'textfield'},
	g1h3 = {header: '发交', dataIndex: 'ccid', editor: 'textfield'},
	g1h4 = {header: '单号', dataIndex: 'id', renderer: Es.renderId},
	g1h5 = {header: '结算日', dataIndex: 'cDate', renderer: Es.renderDate};
var g2h1 = {header: '编号', dataIndex: 'mid', editor: 'textfield'},
	g2h2 = {header: '批号', dataIndex: 'lot', editor: 'textfield'},
	g2h3 = {header: '数量', dataIndex: 'qty', editor: 'numberfield'},
	g2h4 = {header: '备注', dataIndex: 'rem', editor: 'textfield'};
/**
 * 模板主体
 * 主从表模板：主表-明细 支持增删改
 */
//store for grid1, grid2
var s1 = new ex.Store({	autoLoad: true, fields: c1,	url: '/esm/' + m1 + '?' });
var s2 = new ex.Store({	fields: c2,	url: '/esd/' + m2 + '?' });
var s3 = new ex.Store({ fields: c2,	url: '/esd/' + m2 + '?' });	
//grid1
var g1 = Ext.create('ex.Grid', {
	store: s1, columns: [
		new Ext.grid.RowNumberer(),
		g1h1,g1h2,g1h3,g1h4,g1h5,
		{
			xtype:'actioncolumn', width:50,
            items: [{
                iconCls: 'icon-yes', tooltip: '过账', handler: function(grid, rowIndex, colIndex) {
                    var rec = grid.getStore().getAt(rowIndex);
                    alert("过账 " + rec.get('id'));
					rec.set('cDate', Ext.Date.format(new Date(), 'Y-m-d H:i:s'));
					s1.sync();
                }
            },{
                iconCls: 'icon-del', tooltip: '删除', handler: function(grid, rowIndex, colIndex) {
                    Ext.Msg.alert("提示", '确认删除?', function(btn) {
						if (btn=='ok') { s1.removeAt(rowIndex), s1.sync();}
					});				
                }
            }]
		}
	],
	onClick: function(me, rec, index) {
		var id = rec.data.id;
		s2.proxy.url = '' + '/esd/' + m2 + '?id=' + id;
		s3.proxy.url = '' + '/esd/' + m2 + '?id=' + id;
		s2.loadPage(1);
	}
});

//grid2
var g2 = Ext.create('ex.Grid', {
	store: s2, maxHeight: 200, dock: 'bottom', selType: 'rowmodel',
	columns: [
		new Ext.grid.RowNumberer(),
		g2h1,g2h2,g2h3,g2h4
	]
});

//panel
var p = {
	xtype:'panel', title: h1, itemId: t1,	glyph: 99, autoScroll: true, closable: true,
	tbar:[{
		text: "新建", iconCls: 'icon-add', handler: function() {
			s1.insert(0, {mDate: Ext.Date.format(new Date(), 'Y-m-d H:i:s')});
			g1.startEdit(0, 0);			
		}
	},{
		text: '修改', iconCls: 'icon-yes', handler: function () {		
			var rec = g1.getSelectRec();
			if (!rec) return;
			var fm = Ext.create('ex.Form', {
				url: 'esm/' + m1 + '?id=' + rec.data.id,
				layout: 'anchor',padding: 10,
				defaults: {xtype: 'textfield', labelAlign : 'right',labelWidth: 50},
				frame: true, border : false,
				items: [{
						name: 'mDate', fieldLabel:'日期', format: 'Y-m-d', anchor: '100%'
					},{
						name: 'lcid', fieldLabel: '仓库', anchor: '100%'
					},{
						name: 'ccid', fieldLabel: '发交', anchor: '100%'
					},{
						xtype: 'container', layout: 'hbox', defaults: {xtype: 'textfield', labelAlign: 'right', labelWidth: 50},
						items:[
							{name: 'id', fieldLabel: '单号', xtype: 'displayfield', type: 'number', renderer: Es.renderId},
							{name: 'cDate', fieldLabel: '结算日', xtype: 'displayfield'}
						]
					},{
					xtype: 'fieldset', title:'明细', margin: '0 5 5 5', padding: 5,
					defaults: {xtype: 'textfield', labelAlign : 'right'}, 
					items: [,{
						xtype: 'xgrid2',
						store: s3,
						columns: [g2h1,g2h2,g2h3,g2h4]
					}]
				}],
				SubmitCallback: function(){
					s3.sync();
					this.up('window').close();
					s1.load();
					s2.load();
				}
			});
			fm.loadRecord(rec);
			var win = Ext.create('Ext.window.Window', {
				autoScroll: true, title: h1, width: 550, items: [fm]
			}).show();
		}
	},{
		text: '删除', iconCls: 'icon-del', handler: function() {
			var selecRec = g1.getSelectRec();
			if (!selecRec) return;
			Ext.Msg.alert("提示", '确认删除?', 	function(btn) {
				if (btn=='ok') s1.remove(selecRec), s1.sync();
			});	
		}
	},'-',{
		text: '保存', iconCls: 'icon-save', handler: function() {
			s1.sync(), s1.load();
		}
	}],
	items: [g1, g2]
}
return p;