package main

import "io/ioutil"
import "net/http"
import "strconv"
import "log"
import "encoding/json"
import "errors"
import "time"
import "os"

import "encoding/base64"

const htmlData = "PCFkb2N0eXBlIGh0bWw+DQo8aHRtbD4NCg0KPGhlYWQ+DQoJPHRpdGxlPmR0b29scyBVSTwvdGl0bGU+DQoJPG1ldGEgbmFtZT0idmlld3BvcnQiIGNvbnRlbnQ9IndpZHRoPWRldmljZS13aWR0aCI+DQoJPGxpbmsgcmVsPSJzdHlsZXNoZWV0IiBocmVmPSJodHRwczovL25ldGRuYS5ib290c3RyYXBjZG4uY29tL2Jvb3Rzd2F0Y2gvMy4wLjAvc2xhdGUvYm9vdHN0cmFwLm1pbi5jc3MiPg0KCTxzY3JpcHQgdHlwZT0idGV4dC9qYXZhc2NyaXB0IiBzcmM9Imh0dHBzOi8vYWpheC5nb29nbGVhcGlzLmNvbS9hamF4L2xpYnMvanF1ZXJ5LzIuMC4zL2pxdWVyeS5taW4uanMiPjwvc2NyaXB0Pg0KCTxzY3JpcHQgdHlwZT0idGV4dC9qYXZhc2NyaXB0IiBzcmM9Imh0dHBzOi8vbmV0ZG5hLmJvb3RzdHJhcGNkbi5jb20vYm9vdHN0cmFwLzMuMS4xL2pzL2Jvb3RzdHJhcC5taW4uanMiPjwvc2NyaXB0Pg0KCTxzdHlsZSB0eXBlPSJ0ZXh0L2NzcyI+DQoJCWJvZHkgew0KCQkJcGFkZGluZy10b3A6IDIwcHg7DQoJCX0NCgkJLmZvb3RlciB7DQoJCQlib3JkZXItdG9wOiAxcHggc29saWQgI2VlZTsNCgkJCW1hcmdpbi10b3A6IDQwcHg7DQoJCQlwYWRkaW5nLXRvcDogNDBweDsNCgkJCXBhZGRpbmctYm90dG9tOiA0MHB4Ow0KCQl9DQoJCS8qIE1haW4gbWFya2V0aW5nIG1lc3NhZ2UgYW5kIHNpZ24gdXAgYnV0dG9uICovDQoJCQ0KCQkuanVtYm90cm9uIHsNCgkJCXRleHQtYWxpZ246IGNlbnRlcjsNCgkJCWJhY2tncm91bmQtY29sb3I6IHRyYW5zcGFyZW50Ow0KCQl9DQoJCS5qdW1ib3Ryb24gLmJ0biB7DQoJCQlmb250LXNpemU6IDIxcHg7DQoJCQlwYWRkaW5nOiAxNHB4IDI0cHg7DQoJCX0NCgkJLyogQ3VzdG9taXplIHRoZSBuYXYtanVzdGlmaWVkIGxpbmtzIHRvIGJlIGZpbGwgdGhlIGVudGlyZSBzcGFjZSBvZiB0aGUgLm5hdmJhciAqLw0KCQkNCgkJLm5hdi1qdXN0aWZpZWQgew0KCQkJYmFja2dyb3VuZC1jb2xvcjogI2VlZTsNCgkJCWJvcmRlci1yYWRpdXM6IDVweDsNCgkJCWJvcmRlcjogMXB4IHNvbGlkICNjY2M7DQoJCX0NCgkJLm5hdi1qdXN0aWZpZWQgPiBsaSA+IGEgew0KCQkJcGFkZGluZy10b3A6IDE1cHg7DQoJCQlwYWRkaW5nLWJvdHRvbTogMTVweDsNCgkJCWNvbG9yOiAjNzc3Ow0KCQkJZm9udC13ZWlnaHQ6IGJvbGQ7DQoJCQl0ZXh0LWFsaWduOiBjZW50ZXI7DQoJCQlib3JkZXItYm90dG9tOiAxcHggc29saWQgI2Q1ZDVkNTsNCgkJCWJhY2tncm91bmQtY29sb3I6ICNlNWU1ZTU7DQoJCQkvKiBPbGQgYnJvd3NlcnMgKi8NCgkJCQ0KCQkJYmFja2dyb3VuZC1yZXBlYXQ6IHJlcGVhdC14Ow0KCQkJLyogUmVwZWF0IHRoZSBncmFkaWVudCAqLw0KCQkJDQoJCQliYWNrZ3JvdW5kLWltYWdlOiAtbW96LWxpbmVhci1ncmFkaWVudCh0b3AsICNmNWY1ZjUgMCUsICNlNWU1ZTUgMTAwJSk7DQoJCQkvKiBGRjMuNisgKi8NCgkJCQ0KCQkJYmFja2dyb3VuZC1pbWFnZTogLXdlYmtpdC1ncmFkaWVudChsaW5lYXIsIGxlZnQgdG9wLCBsZWZ0IGJvdHRvbSwgY29sb3Itc3RvcCgwJSwgI2Y1ZjVmNSksIGNvbG9yLXN0b3AoMTAwJSwgI2U1ZTVlNSkpOw0KCQkJLyogQ2hyb21lLFNhZmFyaTQrICovDQoJCQkNCgkJCWJhY2tncm91bmQtaW1hZ2U6IC13ZWJraXQtbGluZWFyLWdyYWRpZW50KHRvcCwgI2Y1ZjVmNSAwJSwgI2U1ZTVlNSAxMDAlKTsNCgkJCS8qIENocm9tZSAxMCssU2FmYXJpIDUuMSsgKi8NCgkJCQ0KCQkJYmFja2dyb3VuZC1pbWFnZTogLW1zLWxpbmVhci1ncmFkaWVudCh0b3AsICNmNWY1ZjUgMCUsICNlNWU1ZTUgMTAwJSk7DQoJCQkvKiBJRTEwKyAqLw0KCQkJDQoJCQliYWNrZ3JvdW5kLWltYWdlOiAtby1saW5lYXItZ3JhZGllbnQodG9wLCAjZjVmNWY1IDAlLCAjZTVlNWU1IDEwMCUpOw0KCQkJLyogT3BlcmEgMTEuMTArICovDQoJCQkNCgkJCWZpbHRlcjogcHJvZ2lkOiBEWEltYWdlVHJhbnNmb3JtLk1pY3Jvc29mdC5ncmFkaWVudChzdGFydENvbG9yc3RyPScjZjVmNWY1JywgZW5kQ29sb3JzdHI9JyNlNWU1ZTUnLCBHcmFkaWVudFR5cGU9MCk7DQoJCQkvKiBJRTYtOSAqLw0KCQkJDQoJCQliYWNrZ3JvdW5kLWltYWdlOiBsaW5lYXItZ3JhZGllbnQodG9wLCAjZjVmNWY1IDAlLCAjZTVlNWU1IDEwMCUpOw0KCQkJLyogVzNDICovDQoJCX0NCgkJLm5hdi1qdXN0aWZpZWQgPiAuYWN0aXZlID4gYSwNCgkJLm5hdi1qdXN0aWZpZWQgPiAuYWN0aXZlID4gYTpob3ZlciwNCgkJLm5hdi1qdXN0aWZpZWQgPiAuYWN0aXZlID4gYTpmb2N1cyB7DQoJCQliYWNrZ3JvdW5kLWNvbG9yOiAjZGRkOw0KCQkJYmFja2dyb3VuZC1pbWFnZTogbm9uZTsNCgkJCWJveC1zaGFkb3c6IGluc2V0IDAgM3B4IDdweCByZ2JhKDAsIDAsIDAsIC4xNSk7DQoJCX0NCgkJLm5hdi1qdXN0aWZpZWQgPiBsaTpmaXJzdC1jaGlsZCA+IGEgew0KCQkJYm9yZGVyLXJhZGl1czogNXB4IDVweCAwIDA7DQoJCX0NCgkJLm5hdi1qdXN0aWZpZWQgPiBsaTpsYXN0LWNoaWxkID4gYSB7DQoJCQlib3JkZXItYm90dG9tOiAwOw0KCQkJYm9yZGVyLXJhZGl1czogMCAwIDVweCA1cHg7DQoJCX0NCgkJQG1lZGlhKG1pbi13aWR0aDogNzY4cHgpIHsNCgkJCS5uYXYtanVzdGlmaWVkIHsNCgkJCQltYXgtaGVpZ2h0OiA1MnB4Ow0KCQkJfQ0KCQkJLm5hdi1qdXN0aWZpZWQgPiBsaSA+IGEgew0KCQkJCWJvcmRlci1sZWZ0OiAxcHggc29saWQgI2ZmZjsNCgkJCQlib3JkZXItcmlnaHQ6IDFweCBzb2xpZCAjZDVkNWQ1Ow0KCQkJfQ0KCQkJLm5hdi1qdXN0aWZpZWQgPiBsaTpmaXJzdC1jaGlsZCA+IGEgew0KCQkJCWJvcmRlci1sZWZ0OiAwOw0KCQkJCWJvcmRlci1yYWRpdXM6IDVweCAwIDAgNXB4Ow0KCQkJfQ0KCQkJLm5hdi1qdXN0aWZpZWQgPiBsaTpsYXN0LWNoaWxkID4gYSB7DQoJCQkJYm9yZGVyLXJhZGl1czogMCA1cHggNXB4IDA7DQoJCQkJYm9yZGVyLXJpZ2h0OiAwOw0KCQkJfQ0KCQl9DQoJCS8qIFJlc3BvbnNpdmU6IFBvcnRyYWl0IHRhYmxldHMgYW5kIHVwICovDQoJCQ0KCQlAbWVkaWEgc2NyZWVuIGFuZChtaW4td2lkdGg6IDc2OHB4KSB7DQoJCQkvKiBSZW1vdmUgdGhlIHBhZGRpbmcgd2Ugc2V0IGVhcmxpZXIgKi8NCgkJCQ0KCQkJLm1hc3RoZWFkLA0KCQkJLm1hcmtldGluZywNCgkJCS5mb290ZXIgew0KCQkJCXBhZGRpbmctbGVmdDogMDsNCgkJCQlwYWRkaW5nLXJpZ2h0OiAwOw0KCQkJfQ0KCQl9DQoJPC9zdHlsZT4NCgk8c2NyaXB0IHR5cGU9InRleHQvamF2YXNjcmlwdCI+DQoJCXZhciBjZlRpbWUgPSBuZXcgRGF0ZSgpOw0KCQl2YXIgY3VyZGlyID0gIiINCgkJdmFyIGRpc0FsaXZlID0gZmFsc2UNCgkJdmFyIFNob3dNZW51ID0gInNlYXJjaCINCgkJdmFyIGNVcGxvYWRGID0gIiINCg0KCQlmdW5jdGlvbiB1cGRhdGVDRWNob1N0KCkgew0KCQkJdmFyIGNFQ2hvUmVxID0gew0KCQkJCUFkZHJlc3M6ICQoIiNhZGRyZXNzLWlkIikudmFsKCksDQoJCQkJUG9ydDogJCgiI3BvcnQtaWQiKS52YWwoKSwNCgkJCQlTZXJ2ZXJBRV9UaXRsZTogJCgiI2FldGl0bGUtaWQiKS52YWwoKQ0KCQkJfTsNCgkJCSQuYWpheCh7DQoJCQkJdXJsOiAiL2MtZWNobyIsDQoJCQkJdHlwZTogIlBPU1QiLA0KCQkJCWRhdGE6IEpTT04uc3RyaW5naWZ5KGNFQ2hvUmVxKSwNCgkJCQlkYXRhVHlwZTogImpzb24iDQoJCQl9KS5kb25lKGZ1bmN0aW9uKGpzb25EYXRhKSB7DQoJCQkJY29uc29sZS5sb2coanNvbkRhdGEpDQoJCQkJZGlzQWxpdmUgPSBqc29uRGF0YS5Jc0FsaXZlDQoJCQkJdXBkYXRlVWkoKQ0KCQkJfSkNCgkJfQ0KDQoJCWZ1bmN0aW9uIHVwZGF0ZVVpKCkgew0KCQkJaWYgKGRpc0FsaXZlKSB7DQoJCQkJJCgiI3BhY3Mtc3RhdHVzLWlkIikudGV4dCgib2siKQ0KCQkJfSBlbHNlIHsNCgkJCQkkKCIjcGFjcy1zdGF0dXMtaWQiKS50ZXh0KCJubyBjb25uZWN0aW9uIikNCgkJCX0NCgkJCWlmIChkaXNBbGl2ZSAmJiAoU2hvd01lbnUgPT0gInNlYXJjaCIpKSB7DQoJCQkJJCgiI3NlYXJjaC1wYW5lbCIpLnNob3coKQ0KCQkJCSQoIiNzZWFyY2gtZm9vdGVyIikuc2hvdygpDQoJCQkJJCgiI3NlcmZvb3Rlci1pZCIpLnNob3coKQ0KCQkJCSQoIiNzZXJ0YWJsZS1pZCIpLnNob3coKQ0KCQkJfSBlbHNlIHsNCgkJCQkkKCIjc2VhcmNoLXBhbmVsIikuaGlkZSgpOw0KCQkJCSQoIiNzZWFyY2gtZm9vdGVyIikuaGlkZSgpDQoJCQkJJCgiI3NlcmZvb3Rlci1pZCIpLmhpZGUoKQ0KCQkJCSQoIiNzZXJ0YWJsZS1pZCIpLmhpZGUoKQ0KCQkJfQ0KCQkJaWYgKGRpc0FsaXZlICYmIChTaG93TWVudSA9PSAidXBsb2FkIikpIHsNCgkJCQkkKCIjZmlsZXMtdGFiIikuc2hvdygpDQoJCQkJJCgiI3VwbG9hZGZvb3Rlci1pZCIpLnNob3coKQ0KCQkJfSBlbHNlIHsNCgkJCQkkKCIjdXBsb2FkZm9vdGVyLWlkIikuaGlkZSgpDQoJCQkJJCgiI2ZpbGVzLXRhYiIpLmhpZGUoKQ0KCQkJfQ0KCQkJaWYgKGRpc0FsaXZlICYmIChTaG93TWVudSA9PSAiam9icyIpKSB7DQoJCQkJJCgiI2pvYnNsaXN0Zm9vdGVyLWlkIikuc2hvdygpDQoJCQkJJCgiI2pvYnNsaXN0Iikuc2hvdygpDQoJCQl9IGVsc2Ugew0KCQkJCSQoIiNqb2JzbGlzdGZvb3Rlci1pZCIpLmhpZGUoKQ0KCQkJCSQoIiNqb2JzbGlzdCIpLmhpZGUoKQ0KCQkJfQ0KCQl9DQoNCgkJZnVuY3Rpb24gc2VuZENGaW5kKCkgew0KCQkJdXBkYXRlSm9icygpDQoJCQl2YXIgY2ZkYXQgPSB7DQoJCQkJU2VydmVyU2V0OiB7DQoJCQkJCUFkZHJlc3M6ICQoIiNhZGRyZXNzLWlkIikudmFsKCksDQoJCQkJCVBvcnQ6ICQoIiNwb3J0LWlkIikudmFsKCksDQoJCQkJCVNlcnZlckFFX1RpdGxlOiAkKCIjYWV0aXRsZS1pZCIpLnZhbCgpDQoJCQkJfSwNCgkJCQlQYXRpZW50TmFtZTogJCgiI3BhdGllbnQtbmFtZS1pZCIpLnZhbCgpLA0KCQkJCUFjY2Vzc2lvbk51bWJlcjogJCgiI2FjY2Vzc2lvbi1udW1iZXItaWQiKS52YWwoKSwNCgkJCQlQYXRpZW5EYXRlT2ZCaXJ0aDogJCgiI2RhdGUtYmlydGgtaWQiKS52YWwoKSwNCgkJCQlTdHVkeURhdGU6ICQoIiNzdHVkeS1kYXRlLWlkIikudmFsKCkNCgkJCX07DQoJCQkkLmFqYXgoew0KCQkJCXVybDogIi9jLWZpbmQiLA0KCQkJCXR5cGU6ICJQT1NUIiwNCgkJCQlkYXRhOiBKU09OLnN0cmluZ2lmeShjZmRhdCksDQoJCQkJZGF0YVR5cGU6ICJqc29uIg0KCQkJfSkNCgkJfQ0KDQoJCWZ1bmN0aW9uIHVwZGF0ZUNGaW5kU3QoKSB7DQoJCQkkLmFqYXgoew0KCQkJCXVybDogIi9jLWZpbmRkYXQiLA0KCQkJCXR5cGU6ICJQT1NUIiwNCgkJCQlkYXRhOiBKU09OLnN0cmluZ2lmeShjZlRpbWUpLA0KCQkJCWRhdGFUeXBlOiAianNvbiINCgkJCX0pLmRvbmUoZnVuY3Rpb24oanNvbkRhdGEpIHsNCgkJCQlpZiAoanNvbkRhdGEuUmVmcmVzaCkgew0KCQkJCQljZlRpbWUgPSBqc29uRGF0YS5GVGltZQ0KCQkJCQkkKCIjc2VyY2hyZXNsaXN0IikucmVtb3ZlKCkNCgkJCQkJdmFyIGluZXJIdG1sID0gIiINCgkJCQkJaW5lckh0bWwgPSBpbmVySHRtbC5jb25jYXQoJzx0Ym9keSBpZD0ic2VyY2hyZXNsaXN0Ij4nKQ0KCQkJCQlmb3IgKGluZGV4IGluIGpzb25EYXRhLkNmaW5kUmVzKSB7DQoJCQkJCQlhbiA9IGpzb25EYXRhLkNmaW5kUmVzW2luZGV4XS5BY2Nlc3Npb25OdW1iZXINCgkJCQkJCXBkID0ganNvbkRhdGEuQ2ZpbmRSZXNbaW5kZXhdLlBhdGllbkRhdGVPZkJpcnRoDQoJCQkJCQlzZCA9IGpzb25EYXRhLkNmaW5kUmVzW2luZGV4XS5TdHVkeURhdGUNCgkJCQkJCXBuID0ganNvbkRhdGEuQ2ZpbmRSZXNbaW5kZXhdLlBhdGllbnROYW1lDQoJCQkJCQlpbmVySHRtbCA9IGluZXJIdG1sLmNvbmNhdCgnPHRyIGlkPSJhY2Nlc3MnICsgYW4gKyAnIj48dGQ+JyArIGFuICsgJzwvdGQ+PHRkPicgKyBwbiArICc8L3RkPjx0ZD4nICsgcGQgKyAnPC90ZD48dGQ+JyArIHNkICsgJzwvdGQ+PC90cj4nKQ0KCQkJCQl9DQoJCQkJCWluZXJIdG1sID0gaW5lckh0bWwuY29uY2F0KCcgPC90Ym9keT4nKQ0KCQkJCQkkKCIjc2VydGFibGUtaWQiKS5hcHBlbmQoaW5lckh0bWwpDQoJCQkJCWNvbnNvbGUubG9nKGpzb25EYXRhLkNmaW5kUmVzKQ0KCQkJCX0gZWxzZSB7DQoJCQkJCWNvbnNvbGUubG9nKCJubyBuZWVkIHRvIHVwZGF0ZSIpDQoJCQkJfQ0KCQkJfSkNCgkJfQ0KDQoJCWZ1bmN0aW9uIGNoRGlyKGUpIHsNCgkJCXZhciBuRGlyID0gew0KCQkJCU5ldzogZS5pZCwNCgkJCQlDdXJEaXI6IGN1cmRpcg0KCQkJfTsNCgkJCSQuYWpheCh7DQoJCQkJdXJsOiAiL2NoZCIsDQoJCQkJdHlwZTogIlBPU1QiLA0KCQkJCWRhdGFUeXBlOiAianNvbiIsDQoJCQkJZGF0YTogSlNPTi5zdHJpbmdpZnkobkRpcikNCgkJCX0pLmRvbmUoZGlyVXBkYXRlKQ0KCQl9DQoNCgkJZnVuY3Rpb24gZmlyc1VwZGF0ZSgpIHsNCgkJCXZhciBuRGlyID0gew0KCQkJCU5ldzogIi4iLA0KCQkJCUN1ckRpcjogIi4iDQoJCQl9Ow0KCQkJJC5hamF4KHsNCgkJCQl1cmw6ICIvY2hkIiwNCgkJCQl0eXBlOiAiUE9TVCIsDQoJCQkJZGF0YVR5cGU6ICJqc29uIiwNCgkJCQlkYXRhOiBKU09OLnN0cmluZ2lmeShuRGlyKQ0KCQkJfSkuZG9uZShkaXJVcGRhdGUpDQoJCX0NCg0KCQlmdW5jdGlvbiBkaXJVcGRhdGUoanNvbkRhdGEpIHsNCgkJCSQoIiNmaWxlcy1pZCIpLnJlbW92ZSgpDQoJCQljdXJkaXIgPSBqc29uRGF0YS5DdXJEaXINCgkJCWNvbnNvbGUubG9nKGpzb25EYXRhKQ0KCQkJdmFyIGluZXJIdG1sZmlsZXMgPSAiIg0KCQkJaW5lckh0bWxmaWxlcyA9IGluZXJIdG1sZmlsZXMuY29uY2F0KCc8dGJvZHkgaWQ9ImZpbGVzLWlkIj4nKQ0KCQkJaW5lckh0bWxmaWxlcyA9IGluZXJIdG1sZmlsZXMuY29uY2F0KCc8dHI+PHRkPjwvdGQ+JykNCgkJCWluZXJIdG1sZmlsZXMgPSBpbmVySHRtbGZpbGVzLmNvbmNhdCgnPHRkIG9uY2xpY2s9ImNoRGlyKHRoaXMpIiBpZD0iLi4iPjxpbWcgc3JjPSJodHRwOi8vdXBsb2FkLndpa2ltZWRpYS5vcmcvd2lraXBlZGlhL2NvbW1vbnMvZC9kYy9CbHVlX2ZvbGRlcl9zZXRoX3lhc3Ryb3ZfMDEuc3ZnIiB3aWR0aD0iMzAiIGFsdD0ibG9yZW0iPi4uPC90ZD48L3RyPicpDQoJCQlmb3IgKGluZGV4IGluIGpzb25EYXRhLkZpbGVzKSB7DQoJCQkJbm0gPSBqc29uRGF0YS5GaWxlc1tpbmRleF0uTmFtZQ0KCQkJCWRpID0ganNvbkRhdGEuRmlsZXNbaW5kZXhdLklzRGlyDQoJCQkJaWYgKGpzb25EYXRhLkZpbGVzW2luZGV4XS5Jc0Rpcikgew0KCQkJCQlpbmVySHRtbGZpbGVzID0gaW5lckh0bWxmaWxlcy5jb25jYXQoJzx0ciB3aWR0aD0iNSI+PHRkPjxhICBvbmNsaWNrPSJzZW5kQ1N0b3JlKHRoaXMpIiBpZD0iJyArICdmaScgKyBubSArICcvIiBjbGFzcz0iYnRuIHB1bGwtbGVmdCBidG4tc3VjY2VzcyBidG4teHMiPlVwbG9hZDwvYT48L3RkPicpDQoJCQkJCWluZXJIdG1sZmlsZXMgPSBpbmVySHRtbGZpbGVzLmNvbmNhdCgnPHRkIG9uY2xpY2s9ImNoRGlyKHRoaXMpIiAnICsgJ2lkPSInICsgbm0gKyAnIj48aW1nIHNyYz0iaHR0cDovL3VwbG9hZC53aWtpbWVkaWEub3JnL3dpa2lwZWRpYS9jb21tb25zL2QvZGMvQmx1ZV9mb2xkZXJfc2V0aF95YXN0cm92XzAxLnN2ZyIgd2lkdGg9IjMwIiBhbHQ9ImxvcmVtIj4nICsgbm0gKyAnPC90ZD48L3RyPicpDQoJCQkJfSBlbHNlIHsNCgkJCQkJaW5lckh0bWxmaWxlcyA9IGluZXJIdG1sZmlsZXMuY29uY2F0KCc8dHIgd2lkdGg9IjUiPjx0ZD48YSAgb25jbGljaz0ic2VuZENTdG9yZSh0aGlzKSIgaWQ9IicgKyAnZmknICsgbm0gKyAnIiBjbGFzcz0iYnRuIHB1bGwtbGVmdCBidG4tc3VjY2VzcyBidG4teHMiPlVwbG9hZDwvYT48L3RkPicpDQoJCQkJCWluZXJIdG1sZmlsZXMgPSBpbmVySHRtbGZpbGVzLmNvbmNhdCgnPHRkIG9uY2xpY2s9ImNoRGlyKHRoaXMpIiAnICsgJ2lkPSInICsgbm0gKyAnIj48aW1nIHNyYz0iaHR0cDovL3d3dy5mcmVlY2Fkd2ViLm9yZy93aWtpL2ltYWdlcy8yLzI5L0RvY3VtZW50LW5ldy5zdmciIHdpZHRoPSIzMCIgYWx0PSJsb3JlbSI+JyArIG5tICsgJzwvdGQ+PC90cj4nKQ0KCQkJCX0NCgkJCX0NCgkJCWluZXJIdG1sZmlsZXMgPSBpbmVySHRtbGZpbGVzLmNvbmNhdCgnPC90Ym9keT4nKQ0KCQkJJCgiI2ZpbGVzLXRhYiIpLmFwcGVuZChpbmVySHRtbGZpbGVzKQ0KCQl9DQoNCgkJZnVuY3Rpb24gdXBkYXRlSm9icygpIHsNCgkJCSQoIiNqb2JzbGlzdCIpLmh0bWwoIiIpDQoJCQkkLmFqYXgoew0KCQkJCXVybDogIi9qb2JzIiwNCgkJCQl0eXBlOiAiUE9TVCIsDQoJCQkJZGF0YVR5cGU6ICJqc29uIg0KCQkJfSkuZG9uZShmdW5jdGlvbihqc29uRGF0YSkgew0KCQkJCXZhciBpbmVySHRtbGpvYnMgPSAiIg0KCQkJCWZvciAoaW5kZXggaW4ganNvbkRhdGEpIHsNCgkJCQkJaW5lckh0bWxqb2JzID0gaW5lckh0bWxqb2JzLmNvbmNhdCgnPGxpIGNsYXNzPSJsaXN0LWdyb3VwLWl0ZW0iPicgKyBqc29uRGF0YVtpbmRleF0gKyAnPC9saT4nKQ0KCQkJCX0NCgkJCQkkKCIjam9ic2xpc3QiKS5hcHBlbmQoaW5lckh0bWxqb2JzKQ0KCQkJfSkNCgkJfQ0KDQoJCWZ1bmN0aW9uIHNlbmRDU3RvcmUoZSkgew0KCQkJdmFyIGZwID0gY3VyZGlyICsgJy8nICsgZS5pZC5zdWJzdHJpbmcoMikNCgkJCXZhciBjc2RhdCA9IHsNCgkJCQlTZXJ2ZXJTZXQ6IHsNCgkJCQkJQWRkcmVzczogJCgiI2FkZHJlc3MtaWQiKS52YWwoKSwNCgkJCQkJUG9ydDogJCgiI3BvcnQtaWQiKS52YWwoKSwNCgkJCQkJU2VydmVyQUVfVGl0bGU6ICQoIiNhZXRpdGxlLWlkIikudmFsKCkNCgkJCQl9LA0KCQkJCUZpbGU6IGZwLA0KCQkJfQ0KCQkJY29uc29sZS5sb2coY3NkYXQpDQoJCQkkLmFqYXgoew0KCQkJCXVybDogIi9jLWN0b3JlIiwNCgkJCQl0eXBlOiAiUE9TVCIsDQoJCQkJZGF0YTogSlNPTi5zdHJpbmdpZnkoY3NkYXQpLA0KCQkJCWRhdGFUeXBlOiAianNvbiINCgkJCX0pDQoJCX0NCg0KCQlmdW5jdGlvbiBPbkxvYWQoKSB7DQoJCQljZlRpbWUgPSAwLjA7DQoJCQl1cGRhdGVVaSgpDQoJCQlzZXRJbnRlcnZhbCh1cGRhdGVDRWNob1N0LCA3MDApDQoJCQlzZXRJbnRlcnZhbCh1cGRhdGVDRmluZFN0LCA0MDApDQoJCQlzZXRJbnRlcnZhbCh1cGRhdGVKb2JzLCAyMDAwKQ0KCQkJZmlyc1VwZGF0ZSgpDQoJCX0NCg0KCQlmdW5jdGlvbiBTaG93U2VhcmNoKCkgew0KCQkJU2hvd01lbnUgPSAic2VhcmNoIg0KCQkJdXBkYXRlVWkoKQ0KCQl9DQoNCgkJZnVuY3Rpb24gU2hvd1VwbG9hZCgpIHsNCgkJCVNob3dNZW51ID0gInVwbG9hZCINCgkJCXVwZGF0ZVVpKCkNCgkJfQ0KDQoJCWZ1bmN0aW9uIFNob3dKb2JzKCkgew0KCQkJU2hvd01lbnUgPSAiam9icyINCgkJCXVwZGF0ZUpvYnMoKQ0KCQkJdXBkYXRlVWkoKQ0KCQl9DQoJPC9zY3JpcHQ+DQo8L2hlYWQ+DQoNCjxib2R5IG9ubG9hZD0iT25Mb2FkKCkiPg0KCTxkaXYgY2xhc3M9ImNvbnRhaW5lciI+DQoJCTxkaXYgY2xhc3M9IndlbGwiPg0KCQkJPGRpdiBjbGFzcz0ibmF2YmFyIG5hdmJhci1kZWZhdWx0Ij4NCgkJCQk8ZGl2IGNsYXNzPSJjb250YWluZXIiPg0KCQkJCQk8ZGl2IGNsYXNzPSJuYXZiYXItaGVhZGVyIj4NCgkJCQkJCTxidXR0b24gdHlwZT0iYnV0dG9uIiBjbGFzcz0ibmF2YmFyLXRvZ2dsZSIgZGF0YS10b2dnbGU9ImNvbGxhcHNlIiBkYXRhLXRhcmdldD0iLm5hdmJhci1jb2xsYXBzZSI+IDxzcGFuIGNsYXNzPSJzci1vbmx5Ij5Ub2dnbGUgbmF2aWdhdGlvbjwvc3Bhbj48c3BhbiBjbGFzcz0iaWNvbi1iYXIiPjwvc3Bhbj48c3BhbiBjbGFzcz0iaWNvbi1iYXIiPjwvc3Bhbj48c3BhbiBjbGFzcz0iaWNvbi1iYXIiPjwvc3Bhbj4gPC9idXR0b24+DQoJCQkJCTwvZGl2Pg0KCQkJCQk8ZGl2IGNsYXNzPSJjb2xsYXBzZSBuYXZiYXItY29sbGFwc2UiPg0KCQkJCQkJPHVsIGNsYXNzPSJuYXYgbmF2YmFyLW5hdiI+DQoJCQkJCQkJPGxpIG9uY2xpY2s9IlNob3dTZWFyY2goKSI+IDxhPlNlYXJjaDwvYT4gPC9saT4NCgkJCQkJCQk8bGkgb25jbGljaz0iU2hvd1VwbG9hZCgpIj4gPGE+U3R1ZHkgVXBsb2FkPC9hPiA8L2xpPg0KCQkJCQkJCTxsaSBvbmNsaWNrPSJTaG93Sm9icygpIj4gPGE+Sm9iczwvYT4gPC9saT4NCgkJCQkJCTwvdWw+DQoJCQkJCTwvZGl2Pg0KCQkJCTwvZGl2Pg0KCQkJPC9kaXY+DQoJCQk8ZGl2IGNsYXNzPSJwYW5lbC1mb290ZXIiPkRJQ09NIFNlcnZlciBzZXR0aW5ncyA8L2Rpdj4NCgkJCTx0YWJsZSBjbGFzcz0idGFibGUgdGFibGUtYm9yZGVyZWQgdGFibGUtY29uZGVuc2VkIHRhYmxlLWhvdmVyIHRhYmxlLXN0cmlwZWQiPg0KCQkJCTx0Ym9keT4NCgkJCQkJPHRyPg0KCQkJCQkJPHRkPg0KCQkJCQkJCTxkaXYgY2xhc3M9ImZvcm0tZ3JvdXAiPg0KCQkJCQkJCQk8bGFiZWwgY2xhc3M9ImNvbnRyb2wtbGFiZWwiPkRJQ09NIHNlcnZlciBhZGRyZXNzPC9sYWJlbD4NCgkJCQkJCQkJPGRpdiBjbGFzcz0iY29udHJvbHMiPg0KCQkJCQkJCQkJPGlucHV0IHR5cGU9InRleHQiIGNsYXNzPSJmb3JtLWNvbnRyb2wgaW5wdXQtc20iIGlkPSJhZGRyZXNzLWlkIiB2YWx1ZT0iMjEzLjE2NS45NC4xNTgiPiA8L2Rpdj4NCgkJCQkJCQk8L2Rpdj4NCgkJCQkJCTwvdGQ+DQoJCQkJCQk8dGQ+DQoJCQkJCQkJPGRpdiBjbGFzcz0iZm9ybS1ncm91cCI+DQoJCQkJCQkJCTxsYWJlbCBjbGFzcz0iY29udHJvbC1sYWJlbCI+QUUtVGl0bGU8L2xhYmVsPg0KCQkJCQkJCQk8ZGl2IGNsYXNzPSJjb250cm9scyI+DQoJCQkJCQkJCQk8aW5wdXQgdHlwZT0idGV4dCIgY2xhc3M9ImZvcm0tY29udHJvbCBpbnB1dC1zbSIgaWQ9ImFldGl0bGUtaWQiIHZhbHVlPSJBRV9HRVBBQ1MiPg0KCQkJCQkJCQk8L2Rpdj4NCgkJCQkJCQk8L2Rpdj4NCgkJCQkJCTwvdGQ+DQoJCQkJCQk8dGQ+DQoJCQkJCQkJPGRpdiBjbGFzcz0iZm9ybS1ncm91cCI+DQoJCQkJCQkJCTxsYWJlbCBjbGFzcz0iY29udHJvbC1sYWJlbCI+UG9ydCBudW1iZXI8L2xhYmVsPg0KCQkJCQkJCQk8ZGl2IGNsYXNzPSJjb250cm9scyI+DQoJCQkJCQkJCQk8aW5wdXQgdHlwZT0idGV4dCIgY2xhc3M9ImZvcm0tY29udHJvbCBpbnB1dC1zbSIgaWQ9InBvcnQtaWQiIHZhbHVlPSIxMDQiPg0KCQkJCQkJCQk8L2Rpdj4NCgkJCQkJCQk8L2Rpdj4NCgkJCQkJCTwvdGQ+DQoJCQkJCQk8dGQ+DQoJCQkJCQkJPGRpdiBjbGFzcz0iZm9ybS1ncm91cCI+DQoJCQkJCQkJCTxsYWJlbCBjbGFzcz0iY29udHJvbC1sYWJlbCI+RElDT00gcGluZyBzdGF0dXM6PC9sYWJlbD4NCgkJCQkJCQkJPHA+DQoJCQkJCQkJCQk8bGFiZWwgY2xhc3M9ImNvbnRyb2wtbGFiZWwiIGlkPSJwYWNzLXN0YXR1cy1pZCI+T0s8L2xhYmVsPg0KCQkJCQkJCQk8L3A+DQoJCQkJCQkJPC9kaXY+DQoJCQkJCQk8L3RkPg0KCQkJCQk8L3RyPg0KCQkJCTwvdGJvZHk+DQoJCQk8L3RhYmxlPg0KCQkJPGRpdiBjbGFzcz0icGFuZWwtZm9vdGVyIiBpZD0ic2VhcmNoLWZvb3RlciI+U2VhcmNoIFNldHRpbmdzIDwvZGl2Pg0KCQkJPHRhYmxlIGlkPSJzZWFyY2gtcGFuZWwiIGNsYXNzPSJ0YWJsZSB0YWJsZS1ib3JkZXJlZCB0YWJsZS1jb25kZW5zZWQgdGFibGUtaG92ZXIgdGFibGUtc3RyaXBlZCI+DQoJCQkJPHRib2R5Pg0KCQkJCQk8dHI+DQoJCQkJCQk8dGQ+DQoJCQkJCQkJPGRpdiBjbGFzcz0iZm9ybS1ncm91cCI+DQoJCQkJCQkJCTxsYWJlbCBjbGFzcz0iY29udHJvbC1sYWJlbCI+UGF0aWVudCBuYW1lPC9sYWJlbD4NCgkJCQkJCQkJPGRpdiBjbGFzcz0iY29udHJvbHMiPg0KCQkJCQkJCQkJPGlucHV0IHR5cGU9InRleHQiIGNsYXNzPSJmb3JtLWNvbnRyb2wgaW5wdXQtc20iIGlkPSJwYXRpZW50LW5hbWUtaWQiPiA8L2Rpdj4NCgkJCQkJCQk8L2Rpdj4NCgkJCQkJCTwvdGQ+DQoJCQkJCQk8dGQ+DQoJCQkJCQkJPGRpdiBjbGFzcz0iZm9ybS1ncm91cCI+DQoJCQkJCQkJCTxsYWJlbCBjbGFzcz0iY29udHJvbC1sYWJlbCI+QWNjZXNzaW9uIG51bWJlcjwvbGFiZWw+DQoJCQkJCQkJCTxkaXYgY2xhc3M9ImNvbnRyb2xzIj4NCgkJCQkJCQkJCTxpbnB1dCB0eXBlPSJ0ZXh0IiBjbGFzcz0iZm9ybS1jb250cm9sIGlucHV0LXNtIiBpZD0iYWNjZXNzaW9uLW51bWJlci1pZCI+IDwvZGl2Pg0KCQkJCQkJCTwvZGl2Pg0KCQkJCQkJPC90ZD4NCgkJCQkJCTx0ZD4NCgkJCQkJCQk8ZGl2IGNsYXNzPSJmb3JtLWdyb3VwIj4NCgkJCQkJCQkJPGxhYmVsIGNsYXNzPSJjb250cm9sLWxhYmVsIj5EYXRlIG9mIGJpcnRoPC9sYWJlbD4NCgkJCQkJCQkJPGRpdiBjbGFzcz0iY29udHJvbHMiPg0KCQkJCQkJCQkJPGlucHV0IHR5cGU9InRleHQiIGNsYXNzPSJmb3JtLWNvbnRyb2wgaW5wdXQtc20iIGlkPSJkYXRlLWJpcnRoLWlkIj4gPC9kaXY+DQoJCQkJCQkJPC9kaXY+DQoJCQkJCQk8L3RkPg0KCQkJCQkJPHRkPg0KCQkJCQkJCTxkaXYgY2xhc3M9ImZvcm0tZ3JvdXAiPg0KCQkJCQkJCQk8bGFiZWwgY2xhc3M9ImNvbnRyb2wtbGFiZWwiPlN0dWR5IGRhdGU8L2xhYmVsPg0KCQkJCQkJCQk8ZGl2IGNsYXNzPSJjb250cm9scyI+DQoJCQkJCQkJCQk8aW5wdXQgdHlwZT0idGV4dCIgY2xhc3M9ImZvcm0tY29udHJvbCBpbnB1dC1zbSIgaWQ9InN0dWR5LWRhdGUtaWQiPiA8L2Rpdj4NCgkJCQkJCQk8L2Rpdj4NCgkJCQkJCTwvdGQ+DQoJCQkJCQk8dGQ+DQoJCQkJCQkJPGRpdiBjbGFzcz0iZm9ybS1ncm91cCI+DQoJCQkJCQkJCTxsYWJlbCBjbGFzcz0iY29udHJvbC1sYWJlbCI+IDwvbGFiZWw+DQoJCQkJCQkJCTxkaXYgY2xhc3M9ImNvbnRyb2xzIj4gPGEgb25jbGljaz0ic2VuZENGaW5kKCkiIGNsYXNzPSJidG4gcHVsbC1sZWZ0IGJ0bi1pbmZvIj5GIEkgTiBEPC9hPiA8L2Rpdj4NCgkJCQkJCQk8L2Rpdj4NCgkJCQkJCTwvdGQ+DQoJCQkJCTwvdHI+DQoJCQkJPC90Ym9keT4NCgkJCTwvdGFibGU+DQoJCQk8ZGl2IGNsYXNzPSJwYW5lbC1mb290ZXIiIGlkPSJzZXJmb290ZXItaWQiPkMtRmluZCByZXN1bHQ8L2Rpdj4NCgkJCTx0YWJsZSBjbGFzcz0idGFibGUgdGFibGUtYm9yZGVyZWQgdGFibGUtY29uZGVuc2VkIHRhYmxlLWhvdmVyIHRhYmxlLXN0cmlwZWQiIGlkPSJzZXJ0YWJsZS1pZCI+DQoJCQkJPHRoZWFkPg0KCQkJCQk8dHI+DQoJCQkJCQk8dGggc3R5bGU9IndpZHRoOiAyMCU7Ij5BY2Nlc3Npb24gbnVtYmVyPC90aD4NCgkJCQkJCTx0aD5QYXRpZW50IG5hbWU8L3RoPg0KCQkJCQkJPHRoIHN0eWxlPSJ3aWR0aDogMjUlOyI+UGF0aWVudCBkYXRlIG9mIGJpcnRoPC90aD4NCgkJCQkJCTx0aCBzdHlsZT0id2lkdGg6IDE1JTsiPlN0dWR5IGRhdGU8L3RoPg0KCQkJCQk8L3RyPg0KCQkJCTwvdGhlYWQ+DQoJCQkJPHRib2R5IGlkPSJzZXJjaHJlc2xpc3QiPiA8L3Rib2R5Pg0KCQkJPC90YWJsZT4NCgkJCTxkaXYgY2xhc3M9InBhbmVsLWZvb3RlciIgaWQ9InVwbG9hZGZvb3Rlci1pZCI+VXBsb2FkPC9kaXY+DQoJCQk8dGFibGUgaWQ9ImZpbGVzLXRhYiIgY2xhc3M9InRhYmxlIHRhYmxlLWJvcmRlcmVkIHRhYmxlLXN0cmlwZWQgdGFibGUtY29uZGVuc2VkIj4NCgkJCQk8dGhlYWQ+DQoJCQkJCTx0cj4NCgkJCQkJCTx0aCBzdHlsZT0id2lkdGg6IDElOyI+U2VsZWN0PC90aD4NCgkJCQkJCTx0aD5GaWxlIE5hbWU8L3RoPg0KCQkJCQk8L3RyPg0KCQkJCTwvdGhlYWQ+DQoJCQkJPHRib2R5IGlkPSJmaWxlcy1pZCI+DQoJCQkJCTx0cj4NCgkJCQkJCTx0ZD4NCgkJCQkJCQk8ZGl2IGNsYXNzPSJjaGVja2JveCBwdWxsLWxlZnQiPg0KCQkJCQkJCQk8bGFiZWw+DQoJCQkJCQkJCQk8aW5wdXQgdHlwZT0iY2hlY2tib3giIHZhbHVlPSJ0cnVlIj4gPC9sYWJlbD4NCgkJCQkJCQk8L2Rpdj4NCgkJCQkJCTwvdGQ+DQoJCQkJCQk8dGQ+TWljaGFlbDwvdGQ+DQoJCQkJCQk8dGQ+bm88L3RkPg0KCQkJCQk8L3RyPg0KCQkJCQk8dHI+DQoJCQkJCQk8dGQ+DQoJCQkJCQkJPGRpdiBjbGFzcz0iY2hlY2tib3ggcHVsbC1sZWZ0Ij4NCgkJCQkJCQkJPGxhYmVsPg0KCQkJCQkJCQkJPGlucHV0IHR5cGU9ImNoZWNrYm94IiB2YWx1ZT0idHJ1ZSI+IDwvbGFiZWw+DQoJCQkJCQkJPC9kaXY+DQoJCQkJCQk8L3RkPg0KCQkJCQkJPHRkPk1pY2hhZWw8L3RkPg0KCQkJCQkJPHRkPm5vPC90ZD4NCgkJCQkJPC90cj4NCgkJCQkJPHRyPg0KCQkJCQkJPHRkPg0KCQkJCQkJCTxkaXYgY2xhc3M9ImNoZWNrYm94IHB1bGwtbGVmdCI+DQoJCQkJCQkJCTxsYWJlbD4NCgkJCQkJCQkJCTxpbnB1dCB0eXBlPSJjaGVja2JveCIgdmFsdWU9InRydWUiPiA8L2xhYmVsPg0KCQkJCQkJCTwvZGl2Pg0KCQkJCQkJPC90ZD4NCgkJCQkJCTx0ZD5NaWNoYWVsPC90ZD4NCgkJCQkJCTx0ZD5ubzwvdGQ+DQoJCQkJCTwvdHI+DQoJCQkJPC90Ym9keT4NCgkJCTwvdGFibGU+DQoJCQk8ZGl2IGNsYXNzPSJwYW5lbC1mb290ZXIiIGlkPSJqb2JzbGlzdGZvb3Rlci1pZCI+Sm9iczwvZGl2Pg0KCQkJPHVsIGlkPSJqb2JzbGlzdCIgY2xhc3M9Imxpc3QtZ3JvdXAiPg0KCQkJCTxsaSBjbGFzcz0ibGlzdC1ncm91cC1pdGVtIj5GaXJzdCBJdGVtPC9saT4NCgkJCQk8bGkgY2xhc3M9Imxpc3QtZ3JvdXAtaXRlbSI+U2Vjb25kIEl0ZW08L2xpPg0KCQkJCTxsaSBjbGFzcz0ibGlzdC1ncm91cC1pdGVtIj5UaGlyZCBJdGVtPC9saT4NCgkJCTwvdWw+DQoJCTwvZGl2Pg0KCTwvZGl2Pg0KPC9ib2R5Pg0KDQo8L2h0bWw+"

type FindData struct {
	FTime    int
	CfindRes []FindRes
	Refresh  bool
}

//main srv class
type DJsServ struct {
	jbBal  JobBallancer
	dDisp  DDisp
	echSta EchoRes
	fndTm  int
	fRes   []FindRes
}

//start and init srv
func (srv *DJsServ) Start(listenPort int) error {
	srv.jbBal.Init(&srv.dDisp, srv, srv)
	srv.dDisp.dCln.CallerAE_Title = "AE_DTOOLS"
	http.HandleFunc("/c-echo", srv.cEcho)
	http.HandleFunc("/c-find", srv.cFind)
	http.HandleFunc("/c-finddat", srv.cFindData)
	http.HandleFunc("/c-ctore", srv.cStore)
	http.HandleFunc("/index.html", srv.index)
	http.HandleFunc("/chd", srv.chd)
	http.HandleFunc("/jobs", srv.jobs)
	if err := http.ListenAndServe(":"+strconv.Itoa(listenPort), nil); err != nil {
		return errors.New("error: can't start listen http server")
	}
	return nil
}

//serve cEcho responce
func (srv *DJsServ) cEcho(rwr http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()
	bodyData, err := ioutil.ReadAll(req.Body)
	if err != nil {
		strErr := "error: Can't read http body data"
		http.Error(rwr, err.Error(), http.StatusInternalServerError)
		log.Println(strErr)
		return
	}
	var dec EchoReq
	if err := json.Unmarshal(bodyData, &dec); err != nil {
		strErr := "error: can't parse DicomCEchoRequest data"
		http.Error(rwr, err.Error(), http.StatusInternalServerError)
		log.Println(strErr)
		return
	}

	if err := srv.jbBal.PushJob(dec); err != nil {
		log.Printf("error: can't push job")
		http.Error(rwr, err.Error(), http.StatusInternalServerError)
		return

	}

	js, err := json.Marshal(srv.echSta)
	if err != nil {
		log.Printf("error: can't serialize servise state")
		http.Error(rwr, err.Error(), http.StatusInternalServerError)
		return
	}
	rwr.Write(js)
}

//serve cEcho responce
func (srv *DJsServ) cFind(rwr http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()
	bodyData, err := ioutil.ReadAll(req.Body)
	if err != nil {
		strErr := "error: Can't read http body data"
		http.Error(rwr, err.Error(), http.StatusInternalServerError)
		log.Println(strErr)
		return
	}
	var fr FindReq
	if err := json.Unmarshal(bodyData, &fr); err != nil {
		strErr := "error: can't parse cFind data"
		http.Error(rwr, err.Error(), http.StatusInternalServerError)
		log.Println(strErr)
		return
	}

	if err := srv.jbBal.PushJob(fr); err != nil {
		log.Printf("error: can't push job")
		http.Error(rwr, err.Error(), http.StatusInternalServerError)
		return

	}
	//return non error empty data
	rwr.Write([]byte{0})
}

//serve find data responce
func (srv *DJsServ) cFindData(rwr http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()
	bodyData, err := ioutil.ReadAll(req.Body)
	if err != nil {
		strErr := "error: Can't read http body data"
		http.Error(rwr, err.Error(), http.StatusInternalServerError)
		log.Println(strErr)
		return
	}
	var lctim int
	if err := json.Unmarshal(bodyData, &lctim); err != nil {
		strErr := "error: can't parse time data"
		http.Error(rwr, err.Error(), http.StatusInternalServerError)
		log.Println(strErr)
		return
	}
	if lctim != srv.fndTm {
		fdat := FindData{Refresh: true, CfindRes: srv.fRes, FTime: srv.fndTm}
		js, err := json.Marshal(fdat)
		if err != nil {
			log.Printf("error: can't serialize cfind data")
			http.Error(rwr, err.Error(), http.StatusInternalServerError)
			return
		}
		rwr.Write(js)
	} else {
		fdat := FindData{Refresh: false}
		js, err := json.Marshal(fdat)
		if err != nil {
			log.Printf("error: can't serialize cfind data")
			http.Error(rwr, err.Error(), http.StatusInternalServerError)
			return
		}
		rwr.Write(js)
	}

}

//serve main page request
func (srv *DJsServ) index(rwr http.ResponseWriter, req *http.Request) {
	rwr.Header().Set("Content-Type: text/html", "*")

	content, err := ioutil.ReadFile("index.html")
	if err != nil {
		log.Println("warning: start page not found, return included page")
		val, _ := base64.StdEncoding.DecodeString(htmlData)
		rwr.Write(val)
		return
	}
	rwr.Write(content)
}

func (srv *DJsServ) DispatchError(fjb FaJob) error {
	log.Print("info: dispatch error ")
	log.Println(fjb.ErrorData)
	return nil
}

func (srv *DJsServ) DispatchSuccess(cjb CompJob) error {
	log.Print("info: dispatch success")
	log.Println(cjb)
	switch result := cjb.ResultData.(type) {
	case EchoRes:
		return srv.onCEchoDone(result)
	case []FindRes:
		return srv.onCFindDone(result)
	default:
		log.Printf("info: unexpected job type %v", result)
	}
	return nil
}

func (srv *DJsServ) onCEchoDone(eres EchoRes) error {
	srv.echSta = eres
	return nil
}

func (srv *DJsServ) onCFindDone(fres []FindRes) error {
	srv.fRes = fres
	srv.fndTm = time.Now().Nanosecond()
	return nil
}

func (srv *DJsServ) jobs(rwr http.ResponseWriter, req *http.Request) {
	if jobs, err := srv.jbBal.GetJobsList(); err != nil {
		log.Printf("error: can't get jobs list data")
		http.Error(rwr, err.Error(), http.StatusInternalServerError)
	} else {
		js, err := json.Marshal(jobs)
		if err != nil {
			log.Printf("error: can't serialize jobs list data")
			http.Error(rwr, err.Error(), http.StatusInternalServerError)
			return
		}
		rwr.Write(js)
	}

}
func (srv *DJsServ) chd(rwr http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()

	bodyData, err := ioutil.ReadAll(req.Body)
	if err != nil {
		strErr := "error: Can't read http body data"
		http.Error(rwr, err.Error(), http.StatusInternalServerError)
		log.Println(strErr)
		return
	}
	var chd struct {
		New    string
		CurDir string
	}
	if err := json.Unmarshal(bodyData, &chd); err != nil {
		strErr := "error: can't parse new dir data"
		http.Error(rwr, err.Error(), http.StatusInternalServerError)
		log.Println(strErr)
		return
	}

	dir, ls, err := Lsd(chd.CurDir + string(os.PathSeparator) + chd.New)
	if err != nil {
		http.Error(rwr, err.Error(), http.StatusInternalServerError)
	}
	var rd struct {
		Files  []Finfo
		CurDir string
	}
	rd.CurDir = dir
	rd.Files = ls
	js, err := json.Marshal(rd)
	if err != nil {
		http.Error(rwr, err.Error(), http.StatusInternalServerError)
		return
	}
	rwr.Write(js)
}

func (srv *DJsServ) cStore(rwr http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()
	bodyData, err := ioutil.ReadAll(req.Body)
	if err != nil {
		strErr := "error: Can't read http body data"
		http.Error(rwr, err.Error(), http.StatusInternalServerError)
		log.Println(strErr)
		return
	}
	var cstr CStorReq
	if err := json.Unmarshal(bodyData, &cstr); err != nil {
		strErr := "error: can't parse c-strore date"
		http.Error(rwr, err.Error(), http.StatusInternalServerError)
		log.Println(strErr)
		return
	}
	if err := srv.jbBal.PushJob(cstr); err != nil {
		log.Printf("error: can't push job")
		http.Error(rwr, err.Error(), http.StatusInternalServerError)
		return

	}
	rwr.Write([]byte{0})
}
