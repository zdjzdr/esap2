var Es ={
	me: this,
	iTab: function(){
		return Ext.create('Ext.tab.Panel', {id: 'iCenter', region: 'center'})
	}(),
	iName: null,
	exMod: [
		'ex.Tree', 
		'ex.Grid',
		'ex.Grid2',
		'ex.Store',
		'ex.Form',
	],
	sysLogo: "img/logo.png",
	indexPage: 'tpl/we.html',		
	iScreen: true,
	iMenu: [{
		title: 'ES模块', itemId: 't_es', store: Ext.create('Ext.data.TreeStore', { proxy:{ type:'ajax', url: 'es/menu'}})
	},{	
		title: '人力', itemId: 't_hr'
	},{
		title: '财务', itemId: 't_fi'
	},{	
		title: '销售', itemId: 't_sd'
	},{
		title: '生产', itemId: 't_pp'
	},{
		title: '仓库', itemId: 't_wm'
	},{
		title: '系统', itemId: 't_ss'
	}],
	$: function(itemId){
		return Ext.getCmp(itemId);
	},
	
	loadjs: function(url) {
		var sc = document.createElement("script");
		sc.setAttribute("src",url);
		document.body.appendChild(sc);
	},

	loadTreeStore: function(url) {
		return Ext.create('Ext.data.TreeStore', { proxy:{ type:'ajax', url: url}});
	},

	chkTab: function (rec) {
		return Ext.getCmp('iCenter').getComponent('t' + rec.data.id + '');
	},

	addTab: function (rec) {
		Ext.Ajax.request({
			url: 'tpl/' + rec.data.id + '.js',
			success: function(r) {
				var $code = '(function() {' + r.responseText + '}())';
				var o = eval($code);
				var oo = o.c ? o.c : o;
				Es.iTab.add(oo).show();
			}
		});
	},
	
	aboutMe: function(){Ext.Msg.alert('关于', 'Esap2.0 <br><br>Designed by woylin<br><br><a href="http://iesap.net">--技术支持--</a>');},
	
	chgMd5: function(str) {
		return (str == "") ? "" : hex_md5(str);
	},
	
	chgPwd: function() {
		var fm = Ext.create('Ext.form.Panel', {
			bodyStyle:"padding:6px",
			frame: true,
			jsonSubmit: true,
			url:"es/chgPwd?",			
			defaultType:"textfield",
			defaults: {inputType: 'password'},
			items:[
				{name:"pwd1", fieldLabel: "旧密码"},
				{name:"pwd2", fieldLabel: "新密码", allowBlank: false},
				{name:"pwd3", fieldLabel: "密码确认", allowBlank: false},
				{name:"p1", xtype:'hidden'},
				{name:"p2", xtype:'hidden'},
				{name:"p3", xtype:'hidden'}
			]
		});
		var win =new Ext.Window({
			title:"更改密码",			
			width: 300,
			layout: 'fit',		
			items:[fm],
			buttons:[{
				text:"确定",
				handler:function(){
					var form = fm;
					var form2 = fm.getForm();
					if (form.isValid()) {
						var p1 = form2.findField('p1');
						var p2 = form2.findField('p2');
						var p3 = form2.findField('p3');
						p1.setValue(Es.chgMd5(form2.findField('pwd1').getValue()));
						p2.setValue(Es.chgMd5(form2.findField('pwd2').getValue()));
						p3.setValue(Es.chgMd5(form2.findField('pwd3').getValue()));
						form.submit({
							success: function(form, action) {
								win.close();
								Ext.Msg.alert('提示','更改成功！');
							},
							failure: function(form, action) {
								Ext.Msg.alert('提示',action.result.msg);
							}
						});
					}
				}
			}]
		}).show();
	},
	checklogin: function(){
		Ext.Ajax.request({
			url: '/login',
			success: function(r) {
				var rr =Ext.decode(r.responseText, true);
				if(!rr.success) Es.login();
			},
			failure: function(r) {				
				Es.login();
			}
		});
	},
	login: function() {
		var fm = Ext.create('Ext.form.Panel', {
			bodyStyle:"padding:6px",
			frame: true,
			url:"/login?",			
			defaults:{xtype:"textfield", labelWidth: 50},
			items:[
				{name:"username",fieldLabel: "用户名", allowBlank: false},
				{name:"password", fieldLabel: "密 码", inputType: 'password'},
				{name:"token", xtype:'hidden'},
				{name:"p", xtype:'hidden'},
			]
		});
		var win =new Ext.Window({
			title:"ESAP登陆",
			width: 280,
			layout: 'fit',	
			modal : true,
			closable: false,
			items:[fm],
			buttons:[{
				text:"登陆",
				handler:function(){
					var form = fm;
					var form2 = fm.getForm();
					if (form.isValid()) {
						var p = form2.findField('p');
						p.setValue(Es.chgMd5(form2.findField('password').getValue()));
						form.submit({
							success: function(form, action) {
								win.close();
								Es.iName = Ext.util.Cookies.get("esapUsrDisp")
							},
							failure: function(form, action) {
								Ext.Msg.alert('提示', action.result.msg);
								// win.close();
							}
						});
					}
				}
			}]
		}).show();
	},
	/**
	 * 将Ext.Json.Store对象
	 */  
	store2json: function(store) {
		var arr;
		if (store instanceof Ext.data.Store) {
			arr= new Array();
			store.each(function(rec){
				// alert(rec);
				arr.push(rec.data);
			});
		} else if(store instanceof Array){
			arr = new Array();
			Ext.each(store,function(rec) {
				arr.push(rec.data);
			});
		}
		return Ext.encode(arr);
	},
	
	//通用单号渲染
	renderId: function(value){ return Ext.String.format('No.0000{0}',value) },
	
	//日期渲染
	renderDate: function(v){ return v!='NULL' ? v : '' },
	
	//通用store reader
	reader: {root: 'data', totalProperty: 'total'},	
}