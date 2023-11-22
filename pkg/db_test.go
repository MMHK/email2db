package pkg

import (
	"email2db/tests"
	"fmt"
	"mime"
	"net/url"
	"os"
	"path/filepath"
	"testing"
)

func loadMySQLDSN() string {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	dbname := os.Getenv("DB_DATABASE")
	username := os.Getenv("DB_USERNAME")
	pwd := os.Getenv("DB_PASSWORD")
	tz := os.Getenv("TZ")

	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=%s",
		url.PathEscape(username),
		url.PathEscape(pwd),
		host, port, url.PathEscape(dbname),
		url.QueryEscape(tz))
}

func TestDBHelper_SaveAttachment(t *testing.T) {
	testImage := tests.GetLocalPath("../tests/sample.jpeg")
	fileExt := filepath.Ext(testImage)

	dsn := loadMySQLDSN()

	t.Log(dsn)

	model := &AttachmentModel{
		Name: filepath.Base(testImage),
		Path: fmt.Sprintf("%s%s", MakeUUID(), fileExt),
		MimeType: mime.TypeByExtension(fileExt),
	}

	helper, err := GetDBHelper(&DBConfig{
		&MySQLConfig{
			DSN: dsn,
		},
	})
	if err != nil {
		t.Error(err)
		t.Fail()
		return
	}
	ID, err := helper.SaveAttachment(model)
	if err != nil {
		t.Error(err)
		t.Fail()
		return
	}

	t.Log(ID)
	t.Log("PASS")
}

func TestDBHelper_SaveMail(t *testing.T) {
	dsn := loadMySQLDSN()

	t.Log(dsn)

	model := &MailModel{
		Subject: "RE: [EXTERNAL] [UAT][SpeedyAgency] ZONEä½¿ç”¨å•\u008Fé¡Œ(A4) - #2023011818264365",
		From: "ia.bcc@hk.zurich.com",
		ReplyTo: "ia.bcc@hk.zurich.com",
		To: []string{"jx3isrz6lp+member@uat-ticket.speedyagency.com"},
		Meta: map[string]string{
			"headers": `Received: by mx0107p1las1.sendgrid.net with SMTP id 5hr6D24VHc Thu, 19 Jan 2023 00:52:33 +0000 (UTC)
Received: from m0134868.ppops.net (unknown [148.163.147.111]) by mx0107p1las1.sendgrid.net (Postfix) with ESMTPS id 82268AC12AC for <jx3isrz6lp+member@uat-ticket.speedyagency.com>; Thu, 19 Jan 2023 00:52:33 +0000 (UTC)
Received: from m0134868.ppops.net (m0134868.ppops.net [127.0.0.1]) by pps.tlsfallback (8.17.1.5/8.17.1.5) with ESMTP id 30J0njGV009885 for <jx3isrz6lp+member@uat-ticket.speedyagency.com>; Thu, 19 Jan 2023 11:52:33 +1100
DKIM-Signature: v=1; a=rsa-sha256; c=relaxed/relaxed; d=hk.zurich.com; h=from : to : subject : date : message-id : references : in-reply-to : content-type : mime-version; s=pps201809; bh=W4zkUkGNGVbWGhfuaDBexhVZW/RoW6CBuWNHO80epeA=; b=nCqspphMA5iZ9ovp2EozfFoNRynNbZDzzq11ZdFMYbq5aS6i4c3m9FdsLSgGeVlUB7Ez 7tHsQiFCkRYho20cgft8R6SB8Dr7QfXrTmLbYpvOA5z7HkfQXpTAjFVEEOgtS35O+YYQ 0Uwhpno/nVQ6hQo4qY0QufWg4SwGvFjSaXraDvrDtIBBN7eRr1SCVzTZOECdFfw51dRw qys1a+iBU1YkXgiVMDjM1ojaFMHBia2Veu7b1WBfPwg979MiHfophpV0AYUmYK9R32jX lCU33+m3jQpiGdadOSN76sx1ruYB/vBOxDg9MVCUgm/MdeW+AorkvsQ/PHgR5utPn54A xg== 
Received: from apc01-tyz-obe.outbound.protection.outlook.com (mail-tyzapc01lp2047.outbound.protection.outlook.com [104.47.110.47]) by mx0a-002e5d03.pphosted.com (PPS) with ESMTPS id 3n6cs70t4k-1 (version=TLSv1.2 cipher=ECDHE-RSA-AES256-GCM-SHA384 bits=256 verify=NOT) for <jx3isrz6lp+member@uat-ticket.speedyagency.com>; Thu, 19 Jan 2023 11:52:32 +1100
ARC-Seal: i=1; a=rsa-sha256; s=arcselector9901; d=microsoft.com; cv=none; b=N7TZFFwtMRAIQsvMNJ/kO6IEQwjZSViOdbfhNJ3JOn9vkod2+0xz6ZkwPhKoQVe1Gdh4mMGdJfLtWHAK1f6qQRRS0gOWbAjnBwV0+A1V/DQdrUMQdmsWJ53szoqv8n3QMyCUiAlODlTniQ2CYTCxaBnD4DTg/zLflIAlqkUYofsqsELxDstth/ukIQl07Id0vvcVzm7H/VVvvSF9aT736rXB0hGMv9aKuZLMZRQG1uYvnF+4iUHA/Dc4hT/fJ2wMba19mL9i+YyptbBEQR98aG5salPGF//oYbJoOKWDm7TKvAcwnUsOSisYfjBpHBBu+ApvfAjLwkb8ki+AygwDNQ==
ARC-Message-Signature: i=1; a=rsa-sha256; c=relaxed/relaxed; d=microsoft.com; s=arcselector9901; h=From:Date:Subject:Message-ID:Content-Type:MIME-Version:X-MS-Exchange-AntiSpam-MessageData-ChunkCount:X-MS-Exchange-AntiSpam-MessageData-0:X-MS-Exchange-AntiSpam-MessageData-1; bh=W4zkUkGNGVbWGhfuaDBexhVZW/RoW6CBuWNHO80epeA=; b=BijDxQfZPHAibtCMgaCgtuiSwBUZfUw1OKwIiay2B2Tx9KcOkDMEorJvG7TKP+bGAyVJjlvD6URWGwOEpGIomZgytrmw6PmkCs+Vd05LBaUWJbJjmrrTlXW0NHs8E+pJMaHooFGB8oswH7HTKMJtdtFs1mUEsiIQ2z44HVBbSVDn5fpRm8Qaxq+RZWkoWIf9Hx+yFSHCvD0G3hoIZ0/kC4fCY/uuLco9q7Yn8H1VB10ZrUNO/yNDhNucow97NPwEfe1X9BNNKUQj2ye42VTb0AOqvbuqlooXY9sc80+4mg2VPqwM60eX2F3k5PXvncNS67R/gA5umkGtnCMMyytn7A==
ARC-Authentication-Results: i=1; mx.microsoft.com 1; spf=pass smtp.mailfrom=hk.zurich.com; dmarc=pass action=none header.from=hk.zurich.com; dkim=pass header.d=hk.zurich.com; arc=none
Received: from SG2PR01MB2029.apcprd01.prod.exchangelabs.com (2603:1096:3:21::17) by KL1PR01MB5232.apcprd01.prod.exchangelabs.com (2603:1096:820:d4::14) with Microsoft SMTP Server (version=TLS1_2, cipher=TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384) id 15.20.6002.24; Thu, 19 Jan 2023 00:52:24 +0000
Received: from SG2PR01MB2029.apcprd01.prod.exchangelabs.com ([fe80::620:28:84e2:bcbb]) by SG2PR01MB2029.apcprd01.prod.exchangelabs.com ([fe80::620:28:84e2:bcbb%7]) with mapi id 15.20.6002.025; Thu, 19 Jan 2023 00:52:24 +0000
From: HKZ IABCC <ia.bcc@hk.zurich.com>
To: "jx3isrz6lp+member@uat-ticket.speedyagency.com" <jx3isrz6lp+member@uat-ticket.speedyagency.com>
Subject:  =?big5?B?UkU6IFtFWFRFUk5BTF0gW1VBVF1bU3BlZWR5QWdlbmN5XSBaT05FqM+lzrDdw0Qo?= =?big5?Q?A4)_-_#2023011818264365?=
Thread-Topic:  =?big5?B?W0VYVEVSTkFMXSBbVUFUXVtTcGVlZHlBZ2VuY3ldIFpPTkWoz6XOsN3DRChBNCkg?= =?big5?Q?-_#2023011818264365?=
Thread-Index: AQHZKydiR4RMa9XW30azGypzZw+g7q6k6iWA
Date: Thu, 19 Jan 2023 00:52:23 +0000
Message-ID:  <SG2PR01MB202924E8C17A3557FECE4C0FDBC49@SG2PR01MB2029.apcprd01.prod.exchangelabs.com>
References: <733ac45ea2a257affb59b6f482c51313@speedyagency.demo2.mixmedia.com>
In-Reply-To:  <733ac45ea2a257affb59b6f482c51313@speedyagency.demo2.mixmedia.com>
Accept-Language: en-US
Content-Language: en-US
X-MS-Has-Attach: 
X-MS-TNEF-Correlator: 
msip_labels: MSIP_Label_55c1b5cd-3a3b-4c04-bd57-14cc72f68117_Enabled=true; MSIP_Label_55c1b5cd-3a3b-4c04-bd57-14cc72f68117_SetDate=2023-01-19T00:52:17Z; MSIP_Label_55c1b5cd-3a3b-4c04-bd57-14cc72f68117_Method=Privileged; MSIP_Label_55c1b5cd-3a3b-4c04-bd57-14cc72f68117_Name=55c1b5cd-3a3b-4c04-bd57-14cc72f68117; MSIP_Label_55c1b5cd-3a3b-4c04-bd57-14cc72f68117_SiteId=95d1d810-50cf-4169-8565-6bfba279a0cd; MSIP_Label_55c1b5cd-3a3b-4c04-bd57-14cc72f68117_ActionId=9095258e-5804-46e8-a1ce-349078c1b321; MSIP_Label_55c1b5cd-3a3b-4c04-bd57-14cc72f68117_ContentBits=2`,
			"dkim": `{@hk.zurich.com : pass}`,
			"to": `"jx3isrz6lp+member@uat-ticket.speedyagency.com" <jx3isrz6lp+member@uat-ticket.speedyagency.com>`,
			"html": `<html xmlns:v="urn:schemas-microsoft-com:vml" xmlns:o="urn:schemas-microsoft-com:office:office" xmlns:w="urn:schemas-microsoft-com:office:word" xmlns:m="http://schemas.microsoft.com/office/2004/12/omml" xmlns="http://www.w3.org/TR/REC-html40">
<head>
<meta http-equiv="Content-Type" content="text/html; charset=big5">
<meta name="Generator" content="Microsoft Word 15 (filtered medium)">
<!--[if !mso]><style>v\\:* {behavior:url(#default#VML);}
o\\:* {behavior:url(#default#VML);}
w\\:* {behavior:url(#default#VML);}
.shape {behavior:url(#default#VML);}
</style><![endif]--><style><!--
/* Font Definitions */
@font-face
	{font-family:·s²Ó©úÅé;
	panose-1:2 2 5 0 0 0 0 0 0 0;}
@font-face
	{font-family:"Cambria Math";
	panose-1:2 4 5 3 5 4 6 3 2 4;}
@font-face
	{font-family:Calibri;
	panose-1:2 15 5 2 2 2 4 3 2 4;}
@font-face
	{font-family:"Calibri Light";
	panose-1:2 15 3 2 2 2 4 3 2 4;}
@font-face
	{font-family:"\\@·s²Ó©úÅé";
	panose-1:2 1 6 1 0 1 1 1 1 1;}
/* Style Definitions */
p.MsoNormal, li.MsoNormal, div.MsoNormal
	{margin:0in;
	font-size:11.0pt;
	font-family:"Calibri",sans-serif;}
h3
	{mso-style-priority:9;
	mso-style-link:"Heading 3 Char";
	mso-margin-top-alt:auto;
	margin-right:0in;
	mso-margin-bottom-alt:auto;
	margin-left:0in;
	font-size:13.5pt;
	font-family:"Calibri",sans-serif;
	font-weight:bold;}
a:link, span.MsoHyperlink
	{mso-style-priority:99;
	color:blue;
	text-decoration:underline;}
span.Heading3Char
	{mso-style-name:"Heading 3 Char";
	mso-style-priority:9;
	mso-style-link:"Heading 3";
	font-family:"Calibri Light",sans-serif;
	color:#1F3763;}
span.EmailStyle19
	{mso-style-type:personal-reply;
	font-family:"Calibri",sans-serif;
	color:windowtext;}
p.msipfooter2462c2f4, li.msipfooter2462c2f4, div.msipfooter2462c2f4
	{mso-style-name:msipfooter2462c2f4;
	mso-margin-top-alt:auto;
	margin-right:0in;
	mso-margin-bottom-alt:auto;
	margin-left:0in;
	font-size:12.0pt;
	font-family:"·s²Ó©úÅé",serif;}
.MsoChpDefault
	{mso-style-type:export-only;
	font-size:10.0pt;}
@page WordSection1
	{size:8.5in 11.0in;
	margin:1.0in 1.25in 1.0in 1.25in;}
div.WordSection1
	{page:WordSection1;}
--></style><!--[if gte mso 9]><xml>
<o:shapedefaults v:ext="edit" spidmax="1026" />
</xml><![endif]--><!--[if gte mso 9]><xml>
<o:shapelayout v:ext="edit">
<o:idmap v:ext="edit" data="1" />
</o:shapelayout></xml><![endif]-->
</head>
<body lang="EN-US" link="blue" vlink="purple" style="word-wrap:break-word">
<div class="WordSection1">
<p class="MsoNormal">test<o:p></o:p></p>
<p class="MsoNormal"><o:p>&nbsp;</o:p></p>
<div style="border:none;border-top:solid #E1E1E1 1.0pt;padding:3.0pt 0in 0in 0in">
<p class="MsoNormal"><b>From:</b> SpeedyAgency Member [Phone: 11111111] &lt;jx3isrz6lp+member@uat-ticket.speedyagency.com&gt;
<br>
<b>Sent:</b> Wednesday, January 18, 2023 6:27 PM<br>
<b>To:</b> HKZ IABCC &lt;ia.bcc@hk.zurich.com&gt;<br>
<b>Subject:</b> [EXTERNAL] [UAT][SpeedyAgency] ZONE<span lang="ZH-TW" style="font-family:&quot;·s²Ó©úÅé&quot;,serif">¨Ï¥Î°ÝÃD</span>(A4) - #2023011818264365<o:p></o:p></p>
</div>
<p class="MsoNormal"><o:p>&nbsp;</o:p></p>
<div>
<h3 style="line-height:200%"><span lang="ZH-TW" style="font-family:&quot;·s²Ó©úÅé&quot;,serif">·s¸ÜÃD¤w¸g³Ð«Ø</span><o:p></o:p></h3>
<p style="line-height:200%"><span lang="ZH-TW">Ä¬¾¤¥@¥N²z½s¸¹</span>: sa <o:p></o:p></p>
<p style="line-height:200%"><span lang="ZH-TW">¦WºÙ</span>: echo <o:p></o:p></p>
<p style="line-height:200%"><span lang="ZH-TW">¹q¶l</span>: <a href="mailto:echo@gmail.com">
echo@gmail.com</a> <o:p></o:p></p>
<p style="line-height:200%"><span lang="ZH-TW">¬d¸ßÃþ§O</span>: <span lang="ZH-TW">²{¦æ«O³æ¬ÛÃöªA°È</span> (<span lang="ZH-TW">¦³«O³æ¸¹½X</span>)(A) &gt; ZONE<span lang="ZH-TW">¨Ï¥Î°ÝÃD</span>(A4)
<o:p></o:p></p>
<p style="line-height:200%"><span lang="ZH-TW">«O³æ½s¸¹</span>: ATT1234567ZC <o:p></o:p></p>
<p style="line-height:200%"><span lang="ZH-TW">¬d¸ß¤º®e</span>: <o:p></o:p></p>
<div class="MsoNormal" align="center" style="text-align:center;line-height:200%">
<hr size="2" width="100%" align="center">
</div>
<p style="line-height:200%">TEST<o:p></o:p></p>
</div>
<p class="MsoNormal" style="margin-bottom:12.0pt"><o:p>&nbsp;</o:p></p>
<div style="border:solid #8A8A8A 1.0pt;padding:12.0pt 12.0pt 12.0pt 12.0pt">
<p style="background:#E6E9E8"><span lang="ZH-TW" style="font-size:10.5pt;color:#242447">¬d¸ßÃþ§O</span><span style="font-size:10.5pt;color:#242447">:
</span><span lang="ZH-TW" style="font-size:10.5pt;color:#000066">²{¦æ«O³æ¬ÛÃöªA°È</span><span style="font-size:10.5pt;color:#000066"> (<span lang="ZH-TW">¦³«O³æ¸¹½X</span>)(A) &gt; ZONE<span lang="ZH-TW">¨Ï¥Î°ÝÃD</span>(A4)</span><span style="font-size:10.5pt;color:#242447">
<o:p></o:p></span></p>
<p style="background:#E6E9E8"><span lang="ZH-TW" style="font-size:10.5pt;color:#242447">¥Î¤á¹q¶l</span><span style="font-size:10.5pt;color:#242447">:
<a href="mailto:echo@gmail.com"><span style="color:#000066">echo@gmail.com</span></a>
<o:p></o:p></span></p>
<p style="background:#E6E9E8"><strong><span lang="ZH-TW" style="font-size:10.5pt;font-family:&quot;·s²Ó©úÅé&quot;,serif;color:#242447">½Ð¹Á¸Õ¨Ï¥Î¥H¤U¤èªkÃö³¬¦¹</span></strong><strong><span style="font-size:10.5pt;font-family:&quot;Calibri&quot;,sans-serif;color:#242447"> Ticket</span></strong><strong><span lang="ZH-TW" style="font-size:10.5pt;font-family:&quot;·s²Ó©úÅé&quot;,serif;color:#242447">¡G</span></strong><span style="font-size:10.5pt;color:#242447"><br>
1. <span lang="ZH-TW">¥´¶} </span><a href="https://urldefense.com/v3/__https:/u22595.ct.sendgrid.net/ls/click?upn=7m0RY48hu7JUyF6rq8HOd2zs5qYVk65hfsrgWUI6ofrqAnkaWHpyvg7nRaJMB-2FEk1VpeDtKYqJ8xY3W6Mvx31sk-2FpEZcgqejKh0c1hahJmA-3DUAQS_XdKF7CRdueyCVB2zW-2BavinAJI4HWTBMXy5xtFnPhPUNqyAo3hIJCdEC86v6RhMmeHgezoyRy-2B-2FJy0nJhfNhmsBWKWWBniIanIcu6NUw1gCW7UrV9bRuB9yoFDUsz9vUHyd5MSe76BMy-2FoLo6F0IvLI1irjrtaqnL9NoPML5hCuwe3NyTJK4p9ySDqx-2FFGPco1ht1aVG66r8eq031wDkyng-3D-3D__;!!F1chEv4xxLU1Ecc!j5FKzCT23-3EoFGQrKGLRZQM1ZZ2VJLlHS_pFAeE8SVbJwU_v4B108R344-Zv-f8CgDPlM7E11S5n_wfcymE1l7BoTw_jf60guve0S5g-ik$" target="_blank"><span lang="ZH-TW" style="color:#000066">¦¹Ãì±µ</span><span style="color:#000066">
 https://speedyagency.demo2.mixmedia.com/ticket/jx3isrz6lp/confirm</span></a> <span lang="ZH-TW">
Ãö³¬¦¹</span> Ticket<span lang="ZH-TW">¡C </span><o:p></o:p></span></p>
</div>
<p class="MsoNormal"><img border="0" width="1" height="1" style="width:.0104in;height:.0104in" id="_x0000_i1026" src="https://u22595.ct.sendgrid.net/wf/open?upn=53ryyDMjjPF27OVvz98NSXK4SxjEw4YrdY0-2BfDabTt0mR7axfUeHZecf9hpXcm5hITWILpQpa7R8POruWAzr8QFwGAr2OxMCwpokDpv0yUlU-2BfR13vKylfU8guu0plk9-2Fq1gxuJh9vczlGb7h-2Bh3QArZISzJq5AWsXCFS5uftV5XaZyollWaPHy8xm3Cbvi7AxRitjD5qjdfh4o5AK5FJ9kmS7JAA11hooA11w7LHnY-3D"><o:p></o:p></p>
<p class="MsoNormal"><span style="font-size:12.0pt;font-family:&quot;·s²Ó©úÅé&quot;,serif"><o:p>&nbsp;</o:p></span></p>
<p class="msipfooter2462c2f4" style="margin:0in"><span style="font-size:10.0pt;font-family:&quot;Calibri&quot;,sans-serif;color:black">Confidential \\ Non Personal Data</span><o:p></o:p></p>
</div>
</body>
</html>`,
			"from": `HKZ IABCC <ia.bcc@hk.zurich.com>`,
			"text": `test

From: SpeedyAgency Member [Phone: 11111111] <jx3isrz6lp+member@uat-ticket.speedyagency.com>
Sent: Wednesday, January 18, 2023 6:27 PM
To: HKZ IABCC <ia.bcc@hk.zurich.com>
Subject: [EXTERNAL] [UAT][SpeedyAgency] ZONE¨Ï¥Î°ÝÃD(A4) - #2023011818264365

·s¸ÜÃD¤w¸g³Ð«Ø

Ä¬¾¤¥@¥N²z½s¸¹: sa

¦WºÙ: echo

¹q¶l: echo@gmail.com<mailto:echo@gmail.com>

¬d¸ßÃþ§O: ²{¦æ«O³æ¬ÛÃöªA°È (¦³«O³æ¸¹½X)(A) > ZONE¨Ï¥Î°ÝÃD(A4)

«O³æ½s¸¹: ATT1234567ZC

¬d¸ß¤º®e:

________________________________

TEST


¬d¸ßÃþ§O: ²{¦æ«O³æ¬ÛÃöªA°È (¦³«O³æ¸¹½X)(A) > ZONE¨Ï¥Î°ÝÃD(A4)

¥Î¤á¹q¶l: echo@gmail.com<mailto:echo@gmail.com>

½Ð¹Á¸Õ¨Ï¥Î¥H¤U¤èªkÃö³¬¦¹ Ticket¡G
1. ¥´¶} ¦¹Ãì±µ https://speedyagency.demo2.mixmedia.com/ticket/jx3isrz6lp/confirm<https://urldefense.com/v3/__https:/u22595.ct.sendgrid.net/ls/click?upn=7m0RY48hu7JUyF6rq8HOd2zs5qYVk65hfsrgWUI6ofrqAnkaWHpyvg7nRaJMB-2FEk1VpeDtKYqJ8xY3W6Mvx31sk-2FpEZcgqejKh0c1hahJmA-3DUAQS_XdKF7CRdueyCVB2zW-2BavinAJI4HWTBMXy5xtFnPhPUNqyAo3hIJCdEC86v6RhMmeHgezoyRy-2B-2FJy0nJhfNhmsBWKWWBniIanIcu6NUw1gCW7UrV9bRuB9yoFDUsz9vUHyd5MSe76BMy-2FoLo6F0IvLI1irjrtaqnL9NoPML5hCuwe3NyTJK4p9ySDqx-2FFGPco1ht1aVG66r8eq031wDkyng-3D-3D__;!!F1chEv4xxLU1Ecc!j5FKzCT23-3EoFGQrKGLRZQM1ZZ2VJLlHS_pFAeE8SVbJwU_v4B108R344-Zv-f8CgDPlM7E11S5n_wfcymE1l7BoTw_jf60guve0S5g-ik$> Ãö³¬¦¹ Ticket¡C



Confidential \\ Non Personal Data`,
			"sender_ip": `148.163.147.111`,
			"envelope": `{"to":["jx3isrz6lp+member@uat-ticket.speedyagency.com"],"from":"ia.bcc@hk.zurich.com"}`,
			"attachments": `0`,
			"subject": `RE: [EXTERNAL] [UAT][SpeedyAgency] ZONEä½¿ç”¨å•é¡Œ(A4) - #2023011818264365`,
			"charsets": `{"to":"UTF-8","html":"big5","subject":"UTF-8","from":"UTF-8","text":"big5"}`,
			"SPF": `pass`,
		},
	}

	helper, err := GetDBHelper(&DBConfig{
		&MySQLConfig{
			DSN: dsn,
		},
	})
	if err != nil {
		t.Error(err)
		t.Fail()
		return
	}
	ID, err := helper.SaveMail(model)
	if err != nil {
		t.Error(err)
		t.Fail()
		return
	}

	t.Log(ID)
	t.Log("PASS")
}
