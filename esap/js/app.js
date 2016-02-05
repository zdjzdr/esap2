Ext.require(Es.exMod);
Ext.onReady(function() {
	Ext.QuickTips.init();
	Ext.setGlyphFontFamily();
	Es.checklogin();
	Es.main =Ext.create('Ext.container.Viewport', {
        layout: 'border', padding: 1,
		items: [{
			id: 'iNorth', region: 'north', xtype: 'toolbar', border:0,
			items: ['ESAP2.0-我的工作台','-',{
				xtype: 'textfield', fieldLabel:'搜索', labelAlign: 'right'
				},'->','欢迎登陆 |',{
					xtype: 'button', iconCls: 'icon-user', text: '账号设置',
					arrowAlign: 'right',menu: [{text: '更改密码', handler: Es.chgPwd},{text: '关于', handler: Es.aboutMe}]
				},'-',{
					text: '安全退出', 
					handler:function() {
						Ext.util.Cookies.clear("esapSID");
						Ext.util.Cookies.clear("esapUsrDisp");
						window.location.href="/login";
					}
				}
			]
		},{
			id: 'iWest', region: 'west', xtype: 'panel',title: '导航',
			collapsible: true, split: true, width: 200, //margin: 3,
			layout: {type: 'accordion', titleCollapse: true, animate: true},
			defaults: {xtype: 'xtree'},
			items: Es.iMenu
		}, Es.iTab, {
			id: 'iSouth',region: 'south', xtype:'toolbar', //height: 36, 
			items: [Es.iName,' 已登录到ESAP','->',{
				text:'重新加载',iconCls: 'icon-fresh', handler:function(a){window.location.href="/";}
			},'-',{
				text:'切换全屏',iconCls: 'icon-out', handler:function(a){
					if (Es.iScreen) {
						a.setText('退出全屏');a.setIconCls('icon-in');
						Es.$('iNorth').hide();Es.$('iWest').hide();
					} else {
						a.setText('切换全屏');a.setIconCls('icon-out');
						Es.$('iNorth').show();Es.$('iWest').show();
					}
					Es.iScreen=!Es.iScreen;
				}		
			},'-','Ver 2.0']
		}]
    });
	Es.iTab.add({iconCls: 'icon-house', title: '首页', autoScroll: true, autoLoad: { url: Es.indexPage}});
});