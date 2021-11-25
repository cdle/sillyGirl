// 京东比价
// [rule:raw https?://item\.m\.jd\.[comhk]{2,3}/product/(\d+).html] 
// [rule:raw https?://.+\.jd\.[comhk]{2,3}/(\d+).html] 
// [rule:raw https?://item\.m\.jd\.[comhk]{2,3}/(\d+).html] 
// [rule:raw https?://m\.jingxi\.[comhk]{2,3}/item/jxview\?sku=(\d+)] 
// [rule:raw https?://m\.jingxi\.[comhk]{2,3}.+sku=(\d+)]
// [rule:raw https?://kpl\.m\.jd\.[comhk]{2,3}/product\?wareId=(\d+)]
// [rule:raw https?://wq\.jd\.[comhk]{2,3}/item/view\?sku=(\d+)]
// [rule:raw https?://wqitem\.jd\.[comhk]{2,3}.+sku=(\d+)]
// [rule:raw https?://.+\.jd\.[comhk]{2,3}.+sku=(\d+)]
// [rule:raw https?://.+jd\.[comhk]{2,3}/product/(\d+).html] 
// [rule: jdsku ?]

var _0xod7 = 'jsjiami.com.v6',
    _0xod7_ = ['‮_0xod7'],
    _0xa79c = [_0xod7, 'wrnDlDXDpw==', 'w4nCjx/Dj8OhEsOU', 'w6tnCMOVw4E=', 'UkYDd8KN', 'w6xEA8OTw6A=', 'w45jJcK9Cg==', 'wpxOw4LDryc=', 'HiHDq8OxPg==', 'wodTw4HDmTI=', 'wo3lvbjliKTkubHvvoo=', 'w6jCosKGw7bCqg==', 'NU97D8O9', 'Hg3DiTTCqQ==', 'woEhERNE', 'OUdcAsOJ', 'aytWQD8=', 'cBZbw5TDrTrCsnvDpxs=', 'JCrDtMOJNQ==', 'bU1Ma8OG', 'bcKee8OLXA==', 'wp82wpg=', 'w7AeaDvDqw==', 'Q8KGUsOIfw==', 'wpg3wocMw4M=', 'wrsEwo7CpsO/', 'wovDrsKSHmQ=', 'wp7Dr8KTDXw=', 'woVLw5jDoQ4=', 'AcKiw5U5Cw==', 'OMKnG8KFw6Q=', 'YioJw5jDhA==', 'GDfDr8O7Kg==', 'ZcKMOcKSw6s=', 'wqNNwp5Sw4Y=', 'wqXCucORa8Kk', 'wp/CimdHYQ==', 'woMAe8K2wqM=', 'w4zDoQvDs2A=', 'wp3DtcKTEmcww5DCvz7Chg==', 'w7pJLcKUFsOe', 'HsKceQ==', 'T8O4NXjCuA==', 'B8KhwpHCg34=', 'wqIsw5/DnUA=', 'dRjCmcO+US4i', 'wrEbZMK0wqo=', '5q+X5Lqx57mc5p+/5Lmh5L+55Y6D6IKwJQ==', 'I8K2w5IFJw==', 'cAJdw53Dhg==', 'cEgqR8KQ', 'DMKKFRtA', 'wrPDkw7DvcOE', 'FMOSwr3DncKw', 'w7o/UhHDog==', 'woVTesOKw4I=', 'w7DCnMKOw4g=', 'wo7CvcO2csKZ', 'UcKAdsO6XCA=', 'w5wOwoHCgX9L', 'cMOqa8O5wozDqcKa', 'wrA/OSp6', 'w4PDlMOAwrjCosKSdw==', 'NSLDkCXCjA==', 'dMOga8KUcsOn', 'wqFTw6/DqQo=', 'H3hYD8Oz', 'woEtOhZq', 'QsKNCw==', 'd8K/PcKFw5c=', 'QRMFw67DmA==', 'Z8KVbsOHVQ==', 'Dx/DkMOeIA==', 'w4DDi8OAwpvCoA==', 'w7VEHMOow7Q=', 'w6EOcg==', '5Y+/6aKG5YqI776q', 'A8Kpw5sJCg==', 'wp98wpZxw4g=', 'w4HDhg3DnHY=', 'w5RxGcOkw7Q=', 'LMOWwoPDlg==', 'w4rChSTDmcOpJMOtwq5pwpJP', '6aKR5YmL5Lmf776d', '5Y696LWk5Liy77+F', 'UcO5X8Ozwpo=', 'wp09b8KfwoTCssKlf8O0woo=', 'CzwbwoRcwofDg8OTw5TDlScqaw==', 'ay13', 'LcKwesOgwpnDosK2dcOt', 'UV9vw7sJ', '5YO95aWv6L2M5Ymn5o6b5Luv55id54uf55eG5ou844OZ', 'VcKDUMKqcR9Aw67ChGzCvMOswpbCl8Kp', 'w7INbzzDocOk', '5Yan5Lq85YSx776K', '5Y+P5Y2h5LmF77yv', 'wqTCqcONXcKDYMOXCMKx', 'O1bCm8KTNXDDijzCpCkswpwkwrgJRRjCqXFow44KZsOIwoBzwok8FsKmU8Oswrctw6TCjRrDksOwe8OnZhJYcgU4UsKdGsKhdMOkfMO/YgwPwqZhMMO3VcO3fMOSwqPCisKQwonDvA3DosKIcg==', 'w6fCj3zDpsOmUx4PwrnClxlKfcOBUsKfwrDDsA3Dh8O0wqPDgcOjwofDrsO1wrTCjMKhwpLCj8OpSx3CvE45QsO9dQrCsEnCrn4HCcKHwpXCvsK5V8K8BcKlw4Q9wrMtUixGw6M5w7vCtycpwrXCgQY9JyTCqmgawooEw5DCgMKaecKDURzDqVJ0RUEYM1vCmsOUeMO5wrTCh8KKw5xOwpHDjsOqYcOV', 'V8O6fMKfe8O5IFXDgMOnw5fDr2rCtcOtwoBaN11+XVXDhcOaXWwnwrHCqcOZa8KLw77DoMKIOMObwqPCqgLCocKwH2zCnUsVwrvDhMK+wpZCw4Zww67CkULCq1XDrsOdaEnCinfDu8OFwoxdw6EYUcOawq5xLDjCnSvDiyJbIsK+SGjDtsOAYEbDlMK+DsKbZSpQT8OkSzdvesKrw6lkTUoJVsOKPkU2IFBREcOeWQ==', 'wphZw5PDuG0HwpnCnMOrATE1a8Kpw5vCi8OPwoRXPynDkDU/woRUwpdMw53Cv8ObSlPCggDChgfDmywgw5rCrMORV8KjwpcTw4TDnztAwpTDpzt4TcKrw7HDkCsJw6dhNibDmkMbGk3DrsK0dW53wpwzC114w4MrQgXClwnDhDE1LSU6UsKNacOOw7RgwpvDr8O8w4V1wr7CpBrCrcKXwq5SwrJow7ZtwqUFecOxwqDDt8KdecOhJsORw4HDs8KXw5NIaCYrwp3Cmw==', 'wrw/w4DDl0LDghIC', 'fCFvwpxfwoHCvhVYwqdIwqsWw7ZYXifCpzHCqU/Dm1Y=', 'MUbCisKbMhXClGbCpzUvwoQwwrJGWkHDoHllw4oVfcK5wpJuwod8F8K7GsK4w6J+w6/ChhHChcO0VsOoahtjO1p0E8KoKcKWFMOVX8OlWTcow5RfG8OdDcKjJ8KOw7bDn8KPw4DDrWTDqMOYLsKeeE96wpXClUFGC8K6bsKSIcKXfxbDgsKjX8KewpDDpgPDg1JBw4rDocOQwpDDssKsdRvDtQ9Jwqlnw5fCtcOpQkbCusOowoIYawHCpsKVT8KOLCTCuMO1AUlbwrQowp1Vw67Cn8O2Q8O0a8KZwoYXwp/CgcKuAEfCvMOxdcKuMlTDscOBwoEiwrYyRlEEOcOfbMKywr/CrMKGL3kUw5kuQC/DjFs0w7jDj1wGCjhrAMOrw78+GUXDs2ZuwoHCgMKFRMKNAcKTDAPDhBvDvT3DpBZFwq1SwojCjMORwqZkw5TDvyvCg8KtQsKowoXCkAFMK8K2w43Dr1tnEMOkL1/DvkDDvcKxVl3Cg8KfwpQ7wptWw7/ClcKuRcORISZBNUjDusO4wrYkLcK8XmJjAyDDrcKSH8O/w7HDk8KBLMO9PMK9w6nDvsKmEsK9w78+w67CtA/DrzLChcK7w7twwpluw4fCoTHCrsKhwrRkPk1aacK7W2U6wp0Cw6lswpXCmh3DlcKEcidLWcKiwqzDp8OnLiXCuE9TwodwNMKsw5A0wrXDimLDscORQnpiw5tdb8KowqfDgsKlw4zCvkXCgMKFHsOYOEHDlcOGOzkcw4owdMOrwrtVw7DDjgjDnMKxBw1xKTE6wrXDm8O2wqHDrcOKAznDtnXDgBZ3wqM3U3NYwosAYygPQMKGDcKDwoQZwqHCv8KTwo9kW8Ocw5Qawq97w5zDr8KyJzFSVBHDkMKPw4ggRsKKwpjCvMOAw54Uw7QpKCUpGcKpJMO7OXXDq8O1wpIVw7DDtsKhw4HDlgpdwpbCtMOUwqrCkcK3bsOHLU93w5YQwr4qw4nChBjDksKww5ppwpF1SsOewr0TTcOew75Mw7LCs2XDk8Oywq/CrVFcKGt4VMK/PCXCuA8FD8OHCMK/wq3DiMKDw6XChcKkwp1IDWQ1Dl8ANMOFEsKXOjTCjMKSDMK6wq/CksObw5HCgQFvw4XCjsOacRrDvDUZC8OLw6fDkR8YBELDmcONw7rCsiPCgcKQIMO0wqLDi3/DhcKewqhJQlINw5nCrg3DjR/CpwU1OVTDgcOmQXDDjw0LYMOfKD5Fw4/DusO9UcKww6vCh8OpA8K/wqDDjkpAwpUDw58pA8KZw4YCXcOzbwdOw60aS8OdQwPDlSBtwr3DvcOnEyrDmMOTwqHCgsOdwqDCkm/DgcKIw4rCiUB6wph2MsKiwqLDo8OywrsMOWLDnyBld3TDl8OIIXpaw63CoxDCqEhjw60RTsKHc0XCo191w5/DvsOaH8K9J8OyF8KMw7DDlQTCg8KGIyfDrAU5fzXCkRU6wrzCvxrCv8Kgw6I0dhHDsHB8F0QqOMOXRzMBNwvDviFZwrd2PGnDn0jDuDPDiMK8LcK3BCHDlMOGw4tPw49twpMIMsKqw7hpT3XDkHHCjz1aw4JEYsOtLMKpLsKtf8K5bGBxwq0YwpMew5zCrcOawobCtnjDs8KfwqUsw5XCujLDusOow5HDlMKkRcKZIsK8IVfCmlB5wpnDviTCkzbCkcOLY27CoMOTW8KDc3V3WxAvVUBqwoN0wqM/woVZGn4=', '5p+D5a+d6KOH5L2n6LaUwrfCs2LDp8KXa3k=', '5ZWR5ZKy5ZCi56W3776S', '5puE5pSR5Yye5Y685Lq+5qCl5pa65o+544Os', '5Y645oqC6LWo77yv', 'bifDtzzDlxYgw55f', 'wr1Mwp9/w4HDm8Kfw6ZX', 'bBLDnwLDrw==', 'FcKTDsOVIBLDoRVhJcKl', 'MRB3BMOJV8KxWGxs', 'w6BuwoTCjgjCkldKJ8Ki', '5p+p6aqM5Lu1', '5Y6D5Y2D5Li+', '5q6M5LmI57i95pyv5Liv5L2r5Y6y6IOCw7Y=', 'wqVrYMOrw7Y=', '5Yme5ZOQ5Lqb772N', 'wqMOwoMZw6E=', 'woQvZFHChQ==', 'wqIWwrXChMOW', 'woDCsktcfQ==', 'YDsKwr5f', 'WQZDRjw=', 'wqBXwphDw5c=', 'wpLCkF3CpBd1d8Olw4/DiA==', 'wpLCkF3CpRd0cMOlw4/DgQ==', 'DcKBw70EGg==', 'wr7DmSLDpA==', 'wpg3wpwTw4Y=', 'UcK1wo8TLw==', 'wrg/enfCiVU=', 'IcKUw6wdLHrCpcO5worCuA==', 'w5TDn8OpwpvCjA==', 'dGbDh8Ohw5Y=', 'YkkUW8Ke', '5Y2S5Lik5Lmy77+g', 'w4TDuQwhwrc=', 'wofCg8OPVMKP', 'wrbCm2BkXg==', '6aGv5YmG6Leo77+6', 'Ymx4XcO5', 'w60NcizDq8O7wqs=', 'UHdYU8Ov', 'wo/CoGp0fg==', 'w4w5woHCs2A=', 'woDCrXhFdg==', 'FcKXC8KkNw==', 'w4IQwqrChUU=', 'fjPCicOwSw==', 'wpo8bcKIwrY=', 'wrkueHrCkg==', 'wovDisK+LVo=', 'W8OELUDCrcKewogVwrI=', 'dk0NUsKP', 'wrx7w6LDuiU=', 'bcOhWcOZwo4=', 'UnRzw5wW', 'STs/w7PDvQ==', 'w7jDhAQBwok=', 'UyPCmsOJYw==', 'wqrCvMOXQcKZBcKOTsKoMTPDsmXDshPDs33Cu0UbAcOfF8KKUx9Uw6E0Hm/Dmw==', 'wpokdsKZwok=', 'w61zwqPDisOGe8OXAiw=', 'wqXCo2pMQQ==', 'w6U6ICpqw6g=', 'w4jDsyvDnEI=', 'wrZgwqFkw7M=', 'FMKeCsKyLA4=', 'wowPV0LCtg==', 'wpARwpHCjcO9', 'MMKxwqTCulA=', 'w7pBwpPDq8OD', 'wrAnOjF6', 'elLDrMOUw40=', 'wp8RYMKywqE=', 'LsK0w6kGJg==', 'woEsTsKRwrQ=', 'w5Jtw5tHwo5hBsK+w5c=', 'UsOEfcO4wqg=', 'wqkrwq8Rw5Y=', 'THMcRMKi', 'c8OkWsOCwqI=', 'L8Owwr/DnMKl', 'w6rDoz3Dn3c=', 'woYKw4XDuE4=', 'w6zCjw3DkMOd', 'fD8swqZh', 'w44cw6XDozZPwqzDisOFXzErf8OnwonCnsKTw5IHaWXDk211wrNRwokIw4jCu8OCCxDDiQbDl0zCgXt2wp/Do8KdP8K0wpUYwpPDiyYzw5LCrHh8RcOowq3DgzlKwqg+cg==', 'woPCmlh0Yg==', 'N2vDn8OIw5k=', 'wooMwo41w6M=', 'w6TDrzzDk2s=', 'wpANwp3CicOK', 'w7IPSgvDmw==', 'NsK1BwRI', 'wp9Ww6/DrQ0=', 'Zn1aU8Oh', 'AsKBw5ApDg==', 'w7cGdwnDpw==', 'w6B9L8Oyw74=', 'Qg0aw6rDqQ==', 'I2dQLcOeQw==', 'GxjDqMOVEA==', 'wpRgwoRkw7E=', 'MsONwoDDocKm', 'dgXCksO1', 'VU5iasOs', 'wqTCm8OBR8Ko', 'GsKYIsKgw6Y=', 'wqDCoMOqU8Kt', 'wr43YH/Cgw==', 'VMKlRcO/Qw==', 'bEYjXMKHDsOmGQ==', 'BMKWEMKMLQ==', 'U3oqWsKT', 'DhjDisOCMA==', 'wpQ9w7fDhGA=', 'LhPDsB/Cnw==', '6aKX5Yiw6Laz7760', 'w7MEbirDug==', 'aRLDti/DhQ==', 'wqE2w5nDjFE=', 'dsKAGsKPw6g=', 'H8KQMMK7BA==', 'cDt1WT0=', 'wpFWYsOLw7s=', 'wpoWwqrCksOO', 'TcKXO8Kgw4o=', 'XcOYEWDCpA==', 'dcKQwoIaM8OSdsOD', 'Nm/DhcOew78=', 'w5rDmcOBwr/CpQ==', 'BcKpEcKlw6A=', 'wpsRwrDCjsOe', 'bxImw5rDry8=', 'TyHDjSTDkw==', 'wr4ewqo9w4Q=', 'wq8uGTFl', 'MD/DmCjCrw==', 'wrPCo8OiZsKl', 'w6jDsQbDj24=', 'KsKUw7YpPXE=', 'JQzDhRjCrg==', 'NcKbwrHCj2I=', 'P8Kzw44pBQ==', 'w7U1ZT/DqA==', 'IFLCg8KKMg==', 'wqc4wqo5w4A=', 'wqMMZcKVwqU=', 'wps7d8KUwqY=', 'GcKMwpLCpVE=', 'ERHDqcOCLw==', 'bATClMO4TA==', 'YMKbwpsC', 'LMKGCzhG', 'wogUwowHw48=', 'wpoEe3zCkQ==', 'wp9Uw4TDvjY=', 'wronwosaw6E=', '5Yyb5omD6LWW776w', 'KWnDng==', 'w7XDsDXDnmk=', 'fGnDpg==', 'w4IawpvChg==', 'wr4/w7rDjE4=', 'LQ7DrAU=', 'ehpB', 'VhFgw7fDhA==', 'w7TDrw3Dnmw=', 'wpkLwrbCiA==', 'w5/Ct8OOTyVewoXDs3bDhA==', 'VH3DncO4w5Y=', 'wr4bVMK3woc=', 'w47DgiAhwpQ=', 'w5XCqijDsMO/', 'w7gmQi7Dqw==', 'QAfDnjnDjA==', 'RMKHc8OFfA==', 'w50UwpbCnG8=', 'FMKxGTlF', 'fRTCjg==', 'EsKDwoLCnGk=', 'wpstwoEdw5I=', 'wr4owqksw6I=', 'w5HCmcKmw7fCrg==', 'TcKOFcKN', 'w5LDoMO+wpXCuQ==', 'w7TCnMKIw4jClDPChk7DisKD', 'czTDoA7Diz0lw5JWeg==', 'eMKbwo0RPsOJ', 'wpxOw4TDoS0wwonCkMO+Xg==', 'emRGw4Ml', 'cEgqQg==', 'NBXDuwjCkxpD', 'wr1pw4XDpgM=', 'bwrDrjDDtQ==', 'Qw01woVV', 'EMKTDMKUw5k=', 'wrMxQH7CoQ==', 'Z8KWwowEPg==', 'EXfDscOiw44=', 'bMOVG2fCjg==', 'w4fDowowwqY=', 'jsjianmDi.cfCQom.Cfkv6QIQCXgYY=='];
if (function(_0x3562b5, _0x22dabc, _0x54d52d) {
        function _0x27bd3d(_0x314543, _0x350466, _0x46622b, _0x151a9e, _0x2c3569, _0x512970) {
            _0x350466 = _0x350466 >> 0x8, _0x2c3569 = 'po';
            var _0x6c486c = 'shift',
                _0x3a4bca = 'push',
                _0x512970 = '‮';
            if (_0x350466 < _0x314543) {
                while (--_0x314543) { _0x151a9e = _0x3562b5[_0x6c486c](); if (_0x350466 === _0x314543 && _0x512970 === '‮' && _0x512970['length'] === 0x1) { _0x350466 = _0x151a9e, _0x46622b = _0x3562b5[_0x2c3569 + 'p'](); } else if (_0x350466 && _0x46622b['replace'](/[nDfCQCfkQIQCXgYY=]/g, '') === _0x350466) { _0x3562b5[_0x3a4bca](_0x151a9e); } }
                _0x3562b5[_0x3a4bca](_0x3562b5[_0x6c486c]());
            }
            return 0xb9845;
        };
        return _0x27bd3d(++_0x22dabc, _0x54d52d) >> _0x22dabc ^ _0x54d52d;
    }(_0xa79c, 0x180, 0x18000), _0xa79c) { _0xod7_ = _0xa79c['length'] ^ 0x180; };

function _0x1593(_0x2a016c, _0x209ac1) {
    _0x2a016c = ~~'0x' ['concat'](_0x2a016c['slice'](0x1));
    var _0x159db5 = _0xa79c[_0x2a016c];
    if (_0x1593['DGEFjg'] === undefined) {
        (function() {
            var _0x24f9a2 = typeof window !== 'undefined' ? window : typeof process === 'object' && typeof require === 'function' && typeof global === 'object' ? global : this;
            var _0x736823 = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/=';
            _0x24f9a2['atob'] || (_0x24f9a2['atob'] = function(_0x5ea4a6) { var _0xb7a5ba = String(_0x5ea4a6)['replace'](/=+$/, ''); for (var _0xc6a93 = 0x0, _0x3ad21e, _0x26f826, _0x15c769 = 0x0, _0x569ee1 = ''; _0x26f826 = _0xb7a5ba['charAt'](_0x15c769++); ~_0x26f826 && (_0x3ad21e = _0xc6a93 % 0x4 ? _0x3ad21e * 0x40 + _0x26f826 : _0x26f826, _0xc6a93++ % 0x4) ? _0x569ee1 += String['fromCharCode'](0xff & _0x3ad21e >> (-0x2 * _0xc6a93 & 0x6)) : 0x0) { _0x26f826 = _0x736823['indexOf'](_0x26f826); } return _0x569ee1; });
        }());

        function _0x227eef(_0x251a2c, _0x209ac1) {
            var _0x37e6e4 = [],
                _0x1dc053 = 0x0,
                _0x15ae8e, _0x766924 = '',
                _0x12a28a = '';
            _0x251a2c = atob(_0x251a2c);
            for (var _0x10a620 = 0x0, _0x34bcc7 = _0x251a2c['length']; _0x10a620 < _0x34bcc7; _0x10a620++) { _0x12a28a += '%' + ('00' + _0x251a2c['charCodeAt'](_0x10a620)['toString'](0x10))['slice'](-0x2); }
            _0x251a2c = decodeURIComponent(_0x12a28a);
            for (var _0x3ee24a = 0x0; _0x3ee24a < 0x100; _0x3ee24a++) { _0x37e6e4[_0x3ee24a] = _0x3ee24a; }
            for (_0x3ee24a = 0x0; _0x3ee24a < 0x100; _0x3ee24a++) {
                _0x1dc053 = (_0x1dc053 + _0x37e6e4[_0x3ee24a] + _0x209ac1['charCodeAt'](_0x3ee24a % _0x209ac1['length'])) % 0x100;
                _0x15ae8e = _0x37e6e4[_0x3ee24a];
                _0x37e6e4[_0x3ee24a] = _0x37e6e4[_0x1dc053];
                _0x37e6e4[_0x1dc053] = _0x15ae8e;
            }
            _0x3ee24a = 0x0;
            _0x1dc053 = 0x0;
            for (var _0x4b5539 = 0x0; _0x4b5539 < _0x251a2c['length']; _0x4b5539++) {
                _0x3ee24a = (_0x3ee24a + 0x1) % 0x100;
                _0x1dc053 = (_0x1dc053 + _0x37e6e4[_0x3ee24a]) % 0x100;
                _0x15ae8e = _0x37e6e4[_0x3ee24a];
                _0x37e6e4[_0x3ee24a] = _0x37e6e4[_0x1dc053];
                _0x37e6e4[_0x1dc053] = _0x15ae8e;
                _0x766924 += String['fromCharCode'](_0x251a2c['charCodeAt'](_0x4b5539) ^ _0x37e6e4[(_0x37e6e4[_0x3ee24a] + _0x37e6e4[_0x1dc053]) % 0x100]);
            }
            return _0x766924;
        }
        _0x1593['JXyjOG'] = _0x227eef;
        _0x1593['NZWBzn'] = {};
        _0x1593['DGEFjg'] = !![];
    }
    var _0x35c712 = _0x1593['NZWBzn'][_0x2a016c];
    if (_0x35c712 === undefined) {
        if (_0x1593['psMzdi'] === undefined) { _0x1593['psMzdi'] = !![]; }
        _0x159db5 = _0x1593['JXyjOG'](_0x159db5, _0x209ac1);
        _0x1593['NZWBzn'][_0x2a016c] = _0x159db5;
    } else { _0x159db5 = _0x35c712; }
    return _0x159db5;
};
var now = 0x0;

function main() {
    var _0x447b15 = { 'frHaN': _0x1593('‫0', 'S14v'), 'sgaJQ': function(_0x38c4e5, _0x454f47) { return _0x38c4e5(_0x454f47); }, 'ZUYmu': function(_0x343a16, _0x230d9c) { return _0x343a16 / _0x230d9c; }, 'CAQZh': function(_0x3ad5bc, _0x1deb06) { return _0x3ad5bc === _0x1deb06; }, 'KpeJS': function(_0x4ef1fe, _0x501d8b) { return _0x4ef1fe <= _0x501d8b; }, 'qkhkw': function(_0x2bcf0f, _0x1f4a87) { return _0x2bcf0f + _0x1f4a87; }, 'fqdeE': 'lPCih', 'EKlee': function(_0x3d5638, _0xe6738e) { return _0x3d5638(_0xe6738e); }, 'shtmq': _0x1593('‫1', 'D[0X'), 'EkEkc': function(_0x2d1b47, _0x238782) { return _0x2d1b47 == _0x238782; }, 'enYOO': 'VkqjL', 'aiQnz': function(_0x33ddd6) { return _0x33ddd6(); }, 'vEbHc': function(_0x527559, _0x2e41c0) { return _0x527559 == _0x2e41c0; }, 'FQCQP': function(_0x222e25, _0x3921e3) { return _0x222e25 + _0x3921e3; }, 'yoTmG': _0x1593('‮2', '$^tX'), 'zSSIl': _0x1593('‫3', 'H]Ev'), 'ZyhQK': function(_0x4466c3, _0x651bf5) { return _0x4466c3 + _0x651bf5; }, 'GPsTp': function(_0x5d8bea, _0x5d31b3, _0x49c7e2) { return _0x5d8bea(_0x5d31b3, _0x49c7e2); }, 'uKSZf': function(_0x594893, _0x3123ec) { return _0x594893 !== _0x3123ec; }, 'cBbAk': _0x1593('‮4', 'Zz^R'), 'fMBRR': 'fanli_vip_secret', 'umHgk': 'uuid', 'PGIvg': function(_0x274870, _0x14e2ca) { return _0x274870 + _0x14e2ca; }, 'onBLc': function(_0xed3a2e, _0x4a4510) { return _0xed3a2e + _0x4a4510; }, 'kxXng': function(_0x4d1c92, _0x3b23a9, _0x568704) { return _0x4d1c92(_0x3b23a9, _0x568704); }, 'PkrHt': _0x1593('‫5', 'Zidc'), 'JFnLa': _0x1593('‫6', 'ZSjm'), 'iZIBr': function(_0x423829, _0x46f128, _0x460167) { return _0x423829(_0x46f128, _0x460167); }, 'GLiYt': _0x1593('‮7', 'H]Ev'), 'OUgRJ': function(_0x158319, _0x1a91ae) { return _0x158319 + _0x1a91ae; }, 'caMqB': function(_0xe81a08, _0x46a802) { return _0xe81a08 + _0x46a802; }, 'sptcK': function(_0x427d2a, _0x2945b7, _0x4c2f92) { return _0x427d2a(_0x2945b7, _0x4c2f92); }, 'IAbyT': _0x1593('‫8', 'Zz^R'), 'eKLOa': '&md5=', 'fEPDA': _0x1593('‮9', 'ZSjm'), 'XZZLl': function(_0x32288c, _0x9edf2, _0x2500cf) { return _0x32288c(_0x9edf2, _0x2500cf); }, 'nxCHd': '暂无数据。', 'gVFtn': function(_0x2db5d5, _0x130b88) { return _0x2db5d5 < _0x130b88; }, 'ApKUu': function(_0x34c969, _0x1a01de) { return _0x34c969 + _0x1a01de; }, 'mDGED': function(_0x46f213, _0x1fa72b) { return _0x46f213(_0x1fa72b); }, 'qkAWO': _0x1593('‮a', 'OUQB'), 'cKdyx': function(_0x298d6a, _0x1db942) { return _0x298d6a + _0x1db942; }, 'PZoow': '去领券：', 'QxcdV': function(_0x14666c, _0x5df019) { return _0x14666c + _0x5df019; }, 'xJCve': function(_0x59e3da, _0x4fd871) { return _0x59e3da - _0x4fd871; }, 'ahOlo': function(_0x37b544, _0x3c8da9) { return _0x37b544 > _0x3c8da9; }, 'qqqoa': function(_0x245c02, _0x36b767) { return _0x245c02 + _0x36b767; }, 'rnQxW': function(_0x48b752, _0x3337aa) { return _0x48b752 + _0x3337aa; }, 'QfFBi': function(_0xd011c6, _0x394ee3) { return _0xd011c6 + _0x394ee3; }, 'BnDPJ': _0x1593('‮b', '&FZC'), 'dLcWb': function(_0x2dab91, _0x1abd5e) { return _0x2dab91 % _0x1abd5e; }, 'hxLkv': function(_0x4a5dfe) { return _0x4a5dfe(); }, 'GSMwB': function(_0x171673, _0x13cab4, _0x48de9c) { return _0x171673(_0x13cab4, _0x48de9c); }, 'RMPFg': _0x1593('‮c', '@Heq'), 'gqrZN': '领券购：', 'pAJTa': function(_0x926de5, _0x4efae2) { return _0x926de5 + _0x4efae2; }, 'jTyLa': '去买买：', 'gUtTl': function(_0x3b65f5, _0x2daf82) { return _0x3b65f5 + _0x2daf82; }, 'DxZKZ': _0x1593('‮d', 'J66I'), 'dzDET': function(_0x50176f, _0x1a9ca6) { return _0x50176f != _0x1a9ca6; }, 'XOfNh': _0x1593('‫e', '1sX^'), 'kZrzT': function(_0xdaec66, _0x381a5a) { return _0xdaec66(_0x381a5a); }, 'hEqHo': _0x1593('‮f', 'O&K2'), 'PKfmE': function(_0x3f9b20) { return _0x3f9b20(); }, 'BtGoa': '傻妞返利插件正版授权用户。', 'xwZuE': function(_0x4cdb15, _0x19fb18) { return _0x4cdb15 !== _0x19fb18; }, 'Nsvvd': 'OSYqF', 'zjjmv': function(_0x378ac5, _0x1a2fb4) { return _0x378ac5(_0x1a2fb4); }, 'wOQoO': _0x1593('‮10', 'jh%K'), 'JCRJM': _0x1593('‮11', 'J66I'), 'TTsFk': function(_0x39856b, _0x3d2dc1) { return _0x39856b(_0x3d2dc1); }, 'LnvzD': 'json', 'zvnyp': 'max-age=0', 'oxPAw': '\x22macOS\x22', 'JEows': _0x1593('‫12', 'G&AV'), 'skoxB': _0x1593('‮13', 'qdeo'), 'aSfKT': 'none', 'DOSFQ': _0x1593('‫14', 'T7Zb'), 'ysXip': 'document', 'rcKSU': _0x1593('‮15', 'H]Ev'), 'PkkSi': _0x1593('‫16', 'jh%K'), 'sjDaO': 'jdprice', 'DpHgG': 'ROjbO', 'wjvQi': _0x1593('‮17', 'Zidc'), 'fhMeT': function(_0x4ea94e, _0x3cdce8, _0x2c43f0) { return _0x4ea94e(_0x3cdce8, _0x2c43f0); }, 'jrnRL': function(_0x201ce6, _0x2705f6) { return _0x201ce6 + _0x2705f6; }, 'DEuDC': function(_0x4936a8, _0x502410) { return _0x4936a8(_0x502410); }, 'dTAdM': 'https://imdraw.com:88/jdprice/', 'fSbvB': function(_0x4e134e, _0xd59ba0) { return _0x4e134e != _0xd59ba0; }, 'KBtin': function(_0x457271, _0x204b59) { return _0x457271 + _0x204b59; }, 'bhIbG': _0x1593('‮18', '$^tX'), 'qJyVP': _0x1593('‮19', 'OUQB'), 'bitZn': function(_0x1af48d, _0x5635a4) { return _0x1af48d + _0x5635a4; }, 'FcAzE': _0x1593('‮1a', '1sX^'), 'gYpSG': function(_0xa23d87, _0x2534a2) { return _0xa23d87 + _0x2534a2; }, 'UobgF': function(_0x516186, _0x5097f9) { return _0x516186 + _0x5097f9; }, 'UwARU': function(_0x5155b8, _0x125b18) { return _0x5155b8 + _0x125b18; }, 'laLrk': function(_0xf75262, _0x554464) { return _0xf75262(_0x554464); }, 'uuPLp': 'store', 'khqkf': 'highest', 'oCZfc': function(_0x56af9d, _0x37b569) { return _0x56af9d * _0x37b569; }, 'vymjT': _0x1593('‫1b', 'EJt%'), 'LgBGw': function(_0x8e584a, _0x28f0db) { return _0x8e584a * _0x28f0db; }, 'UABCs': _0x1593('‮1c', '!yu5'), 'HQiZT': function(_0x5c56ab, _0x4b0fb4) { return _0x5c56ab < _0x4b0fb4; }, 'yBVgL': _0x1593('‮1d', 'EJt%'), 'uYdgf': _0x1593('‫1e', 'WUnd'), 'JXgog': function(_0xbc6571, _0x41d19b) { return _0xbc6571 === _0x41d19b; }, 'TObKV': _0x1593('‫1f', '1sX^'), 'nEdGR': 'price', 'EuCVO': function(_0x10b897, _0x1e1e83) { return _0x10b897 < _0x1e1e83; }, 'yNfJm': _0x1593('‮20', 'T7Zb'), 'sIMKk': function(_0x882715, _0x40cd2b) { return _0x882715 === _0x40cd2b; }, 'UPZKS': _0x1593('‮21', 'Zidc'), 'jbaOb': '最低价', 'AbOpL': function(_0x1cf5fd, _0x47252e) { return _0x1cf5fd(_0x47252e); }, 'owBUG': _0x1593('‮22', 'aE9n'), 'WsIxb': function(_0xd929fe, _0x2e44d2) { return _0xd929fe(_0x2e44d2); }, 'WOVME': 'xnRjE', 'QZITM': function(_0x4810d9, _0x1aad09) { return _0x4810d9 === _0x1aad09; }, 'cQNAz': _0x1593('‮23', 'S14v'), 'CPFGH': function(_0x4db149, _0x4374b1) { return _0x4db149 < _0x4374b1; }, 'ChmqT': function(_0x50a57f, _0x559c56) { return _0x50a57f != _0x559c56; }, 'QUnjA': '618', 'lLaSQ': function(_0x2afdbb, _0x496421) { return _0x2afdbb === _0x496421; }, 'zyGWZ': _0x1593('‮24', 'LlRA'), 'qVTdk': function(_0x546771, _0x1288e2) { return _0x546771 + _0x1288e2; }, 'kojUp': function(_0x24bb0e, _0x518f45) { return _0x24bb0e + _0x518f45; }, 'YZQyq': function(_0xe55fa2, _0x23a00e) { return _0xe55fa2 + _0x23a00e; }, 'HqGbO': function(_0x3e341c, _0x4b603d) { return _0x3e341c + _0x4b603d; }, 'RzKFE': _0x1593('‫25', ']!fw'), 'JNwJh': function(_0x485dac, _0x9a8025) { return _0x485dac !== _0x9a8025; }, 'IWnUy': _0x1593('‫26', 't^J%'), 'iwsmL': _0x1593('‫27', '4]QC'), 'nluzN': 'VOqLX', 'cWxEB': function(_0x4abdc2, _0xbbe170) { return _0x4abdc2 !== _0xbbe170; }, 'YVRFA': 'zjsjw', 'aWXgX': _0x1593('‫28', 'kx##'), 'eGJKn': function(_0x41ead2, _0x378027) { return _0x41ead2 + _0x378027; }, 'wroKh': function(_0x1d2009, _0x2150fc) { return _0x1d2009 + _0x2150fc; }, 'LmSnZ': function(_0x4c38de, _0x31a7ed) { return _0x4c38de + _0x31a7ed; }, 'nVKHU': function(_0x44cd1a, _0x354425) { return _0x44cd1a(_0x354425); } };
    var _0x48cd99 = function(_0x2b00d4) {
        if (_0x447b15['fqdeE'] !== _0x1593('‮29', 'gWcR')) {
            var _0x23e980 = _0x16f16f[_0x8e183][_0x447b15[_0x1593('‮2a', 'H]Ev')]] * 0x3e8;
            var _0x372029 = _0x447b15['sgaJQ'](time, _0x23e980)['split']('\x20')[0x0];
            var _0x291211 = Math['round'](_0x447b15[_0x1593('‫2b', 'ekB8')](_0x16f16f[_0x8e183][_0x1593('‫2c', '!yu5')], 0x64));
            _0x447b15['CAQZh'](_0x372029, _0x1593('‫2d', 'kdA6')) ? _0x303875 = _0x291211 : '';
            _0x372029 === _0x1593('‫2e', 'kdA6') ? _0x199ec1 = _0x291211 : '';
            if (dayDiff(_0x372029) < 0x1f && _0x447b15[_0x1593('‫2f', 'OUQB')](_0x291211, _0xcb265f['price'])) {
                _0xcb265f['price'] = _0x291211;
                _0xcb265f[_0x1593('‫30', 'J66I')] = _0x372029;
            }
        } else {
            if (!_0x447b15['EKlee'](get, _0x447b15[_0x1593('‮31', 't^J%')])) { return ![]; }
            var _0x420ff1 = ![];
            var _0x17da19 = _0x447b15[_0x1593('‮32', '%E57')](get, _0x1593('‫33', '4]QC'));
            var _0x9b6b2b = new Date()[_0x1593('‫34', 'OUQB')]();
            if (_0x447b15['EkEkc'](_0x9b6b2b % 0x7, 0x0) || _0x2b00d4) { if (_0x447b15[_0x1593('‮35', 'ZiSQ')] === _0x447b15[_0x1593('‮36', 'MPXv')]) { _0x17da19 = _0x447b15[_0x1593('‮37', '1sX^')](_0x3fcd44); } else { _0x2ae979 += _0x447b15['qkhkw'](_0x1593('‮38', 'Zidc'), _0x59f556[_0x1593('‮39', 'Z[yY')]); } }
            _0x17da19 = _0x447b15[_0x1593('‮3a', 'O&K2')](parseInt, _0x17da19);
            if (_0x447b15['vEbHc'](_0x17da19 % 0x17, 0x0) && _0x17da19) { _0x420ff1 = !![]; }
            return _0x420ff1;
        }
    };
    var _0x3fcd44 = function() { var _0x4bb61b = { 'mkIwN': function(_0xa145d8, _0x20511e) { return _0x447b15[_0x1593('‫3b', 'gWcR')](_0xa145d8, _0x20511e); }, 'lOppc': _0x1593('‮3c', 'kx##'), 'bEtkb': '去买买：' }; var _0x3e9d51 = _0x447b15[_0x1593('‫3d', 'Y^Ra')](bucketGet, 'qq', _0x1593('‮3e', '@Heq')); if (!_0x3e9d51) { if (_0x447b15[_0x1593('‮3f', 'Y^Ra')](_0x447b15[_0x1593('‮40', 'gWcR')], _0x447b15[_0x1593('‫41', 'AQJA')])) { if (_0x50508e) { _0x2ae979 += _0x4bb61b['mkIwN'](_0x4bb61b[_0x1593('‮42', 'gWcR')], _0x59f556[_0x1593('‮43', '&FZC')]); } else { _0x2ae979 += _0x4bb61b[_0x1593('‫44', 'AQJA')](_0x4bb61b[_0x1593('‮45', '6Qwh')], _0x59f556[_0x1593('‫46', 'Zidc')]); } } else { return 0x0; } } var _0x6778cc = _0x3e9d51[_0x1593('‮47', '4]QC')]('&')[0x0]; var _0xe11326 = get(_0x447b15[_0x1593('‫48', '2bIO')]); var _0xa66dbe = bucketGet(_0x1593('‫49', 'i*ZX'), _0x447b15[_0x1593('‫4a', '1sX^')]); var _0x272685 = _0x447b15[_0x1593('‫4b', 'qdeo')](_0x447b15[_0x1593('‮4c', 'Zz^R')](_0x6778cc, _0xe11326) + _0xa66dbe, _0x447b15[_0x1593('‮4d', 'ZSjm')](call, _0x447b15['PkrHt'], _0x447b15[_0x1593('‫4e', 'R#xJ')])); var _0x52fe58 = _0x447b15['iZIBr'](call, _0x447b15['GLiYt'], _0x272685); var _0x1bc70a = request({ 'url': _0x447b15[_0x1593('‮4f', 'Z[yY')](_0x447b15[_0x1593('‫50', '6Qwh')](_0x447b15['OUgRJ'](_0x447b15['caMqB'](_0x1593('‮51', 'O&K2') + _0x447b15[_0x1593('‮52', 'Zidc')](call, _0x1593('‫53', 't9Qt'), '') + _0x447b15[_0x1593('‮54', 'gWcR')] + _0x6778cc, '&secret='), _0xe11326), _0x1593('‫55', 'f(aS')) + _0xa66dbe + _0x447b15['eKLOa'], _0x52fe58) }); if (_0x1bc70a) { if (_0x1593('‮56', 'kdA6') === _0x447b15[_0x1593('‫57', '!yu5')]) { _0x447b15['XZZLl'](set, _0x1593('‫58', '&FZC'), _0x1bc70a); return +_0x1bc70a; } else { if (_0x50508e) { _0x2ae979 += _0x447b15[_0x1593('‫59', '4]QC')](_0x447b15[_0x1593('‫5a', 'kx##')], _0x59f556[_0x1593('‫5b', 'lqB9')]); } else { _0x2ae979 += _0x447b15[_0x1593('‫5c', 't9Qt')] + _0x59f556[_0x1593('‮5d', 'f(aS')]; } } } return 0x0; };
    var _0x4289e9 = _0x447b15[_0x1593('‮5e', 'MPXv')](param, 0x1);
    if (_0x447b15[_0x1593('‫5f', 'Zidc')](_0x4289e9, _0x447b15[_0x1593('‮60', 'OUQB')])) { if (_0x447b15[_0x1593('‫61', 'Zidc')](GetChatID) == _0x1593('‮62', 't^J%')) { _0x447b15[_0x1593('‫63', 'Zz^R')](Continue); return; } if (_0x447b15['kZrzT'](_0x48cd99, _0x4289e9)) { _0x447b15['kZrzT'](sendText, _0x447b15[_0x1593('‫64', 't^J%')]); } else { if (_0x447b15['xwZuE'](_0x447b15['Nsvvd'], _0x1593('‫65', '1sX^'))) { sendText(_0x447b15['nxCHd']); } else { _0x447b15['zjjmv'](sendText, _0x447b15[_0x1593('‮66', 'Zz^R')]); } } return; }
    var _0x2ae979 = '';
    var _0x29fa33 = _0x447b15['gUtTl'](_0x447b15[_0x1593('‮67', 'S14v')], _0x4289e9) + _0x447b15[_0x1593('‮68', 'kdA6')];
    var _0x3940ec = _0x447b15[_0x1593('‮69', 'T7Zb')](request, { 'url': _0x29fa33, 'dataType': _0x447b15[_0x1593('‫6a', 'D[0X')], 'headers': { 'Connection': 'keep-alive', 'Cache-Control': _0x447b15[_0x1593('‫6b', 'H]Ev')], 'sec-ch-ua': _0x1593('‮6c', 'qdeo'), 'sec-ch-ua-mobile': '?0', 'sec-ch-ua-platform': _0x447b15[_0x1593('‫6d', 'gWcR')], 'Upgrade-Insecure-Requests': '1', 'User-Agent': _0x447b15['JEows'], 'Accept': _0x447b15[_0x1593('‮6e', '663P')], 'Sec-Fetch-Site': _0x447b15[_0x1593('‮6f', 't^J%')], 'Sec-Fetch-Mode': _0x447b15[_0x1593('‫70', 'kdA6')], 'Sec-Fetch-User': '?1', 'Sec-Fetch-Dest': _0x447b15[_0x1593('‫71', 'kx##')], 'Accept-Language': _0x447b15[_0x1593('‫72', '@Heq')], 'Cookie': _0x447b15[_0x1593('‫73', ']!fw')] } });
    var _0x233295 = _0x447b15['TTsFk'](cancall, _0x447b15[_0x1593('‮74', 'qdeo')]);
    var _0x59f556 = undefined;
    var _0x50508e = ![];
    _0x48cd99 = _0x48cd99();
    var _0x2cf814 = ![];
    if (_0x48cd99) { _0x2cf814 = !![]; if (!_0x233295) { if (_0x447b15[_0x1593('‮75', 'Y^Ra')](_0x447b15['DpHgG'], _0x447b15[_0x1593('‮76', 'OUQB')])) { sendText(_0x447b15[_0x1593('‫77', '@Heq')]); } else { var _0x5da47a = ''; for (var _0x14834d = 0x0; _0x447b15[_0x1593('‫78', ']]PH')](_0x14834d, len - _0x447b15[_0x1593('‮79', 'R#xJ')](str, '')[_0x1593('‫7a', 'aE9n')]); _0x14834d++) { _0x5da47a += '\x20'; } return _0x5da47a; } } else { var _0x35cfb4 = _0x447b15[_0x1593('‫7b', 'q$H%')](call, 'jdprice', _0x4289e9); if (_0x35cfb4) { _0x59f556 = eval(_0x447b15['jrnRL']('(', _0x35cfb4) + ')'); } } }
    if (!_0x59f556) {
        _0x2cf814 = ![];
        _0x59f556 = _0x447b15[_0x1593('‫7c', '!yu5')](request, { 'url': _0x447b15[_0x1593('‫7d', 'S14v')](_0x447b15['dTAdM'], _0x4289e9), 'dataType': _0x1593('‮7e', '6Qwh') });
    }
    if (_0x59f556 && _0x59f556[_0x1593('‮7f', 'Y^Ra')]) {
        if (_0x447b15[_0x1593('‫80', 'O&K2')](_0x59f556[_0x1593('‮81', '@JFW')], _0x59f556['final'])) { _0x50508e = !![]; }
        _0x2ae979 += _0x447b15['KBtin'](_0x447b15[_0x1593('‫82', 'O&K2')], _0x59f556[_0x1593('‮83', '4]QC')]) + '\x0a\x0a';
        now = _0x59f556['price'];
    }
    if (!_0x3940ec) {
        _0x2ae979 += _0x447b15[_0x1593('‫84', 'WUnd')];
        if (_0x59f556 && _0x59f556[_0x1593('‮85', '1sX^')]) {
            if (_0x48cd99) { if (_0x2cf814) { if (_0x50508e) { _0x2ae979 += _0x447b15[_0x1593('‫86', '&FZC')](_0x447b15[_0x1593('‫87', '1sX^')], _0x59f556[_0x1593('‮88', 'q$H%')]); } else { _0x2ae979 += _0x447b15[_0x1593('‫89', 'T7Zb')] + _0x59f556[_0x1593('‮8a', 'pRJH')]; } } else { if (_0x50508e) { _0x2ae979 += _0x1593('‫8b', 'MPXv') + _0x59f556[_0x1593('‫8c', '@Heq')]; } else { _0x2ae979 += _0x447b15['gYpSG'](_0x447b15[_0x1593('‮8d', 'EJt%')], _0x59f556[_0x1593('‮8e', 'T7Zb')]); } } } else { if (_0x50508e) { _0x2ae979 += _0x447b15[_0x1593('‫8f', '$Ps&')](_0x447b15[_0x1593('‮90', '&FZC')], _0x59f556[_0x1593('‮91', 'ekB8')]); } else { _0x2ae979 += _0x447b15['UwARU'](_0x447b15[_0x1593('‫92', 'LlRA')], _0x59f556[_0x1593('‮93', 'kx##')]); } }
            _0x447b15['laLrk'](sendText, _0x2ae979);
        } else { _0x447b15['laLrk'](sendText, _0x447b15[_0x1593('‫94', '$Ps&')]); }
        return;
    }
    var _0x303875 = 0x0,
        _0x199ec1 = 0x0;
    var _0xcb265f = { 'price': 0x5f5e0ff, 'text': '' };
    var _0x16f16f = _0x3940ec['promo'];
    var _0x22cfbf = _0x3940ec[_0x447b15[_0x1593('‮95', 'i*ZX')]][0x1];
    var _0x5e4a3f = _0x3940ec[_0x1593('‮96', '%E57')];
    var _0x2f79a1 = {};
    if (_0x22cfbf) { _0x2f79a1 = { 'max': Math[_0x1593('‫97', '663P')](_0x22cfbf[_0x447b15[_0x1593('‮98', 'ZiSQ')]]), 'maxt': time(_0x447b15[_0x1593('‮99', '@JFW')](_0x22cfbf[_0x447b15['vymjT']], 0x3e8)), 'min': Math[_0x1593('‫9a', 'kx##')](_0x22cfbf[_0x1593('‫9b', 'R#xJ')]), 'mint': time(_0x447b15[_0x1593('‮9c', 'EJt%')](parseInt(_0x22cfbf[_0x447b15[_0x1593('‮9d', 't^J%')]]), 0x3e8)) }; }
    if (!now && _0x22cfbf) { if (_0x447b15['CAQZh']('PuoJl', 'PuoJl')) { now = _0x447b15[_0x1593('‫9e', 'f(aS')](parseFloat, _0x22cfbf['current_price']); } else { _0x447b15[_0x1593('‮9f', 'pRJH')](sendText, _0x447b15[_0x1593('‫a0', 'O&K2')]); } }
    if (_0x16f16f) {
        for (var _0x8e183 = 0x0; _0x447b15[_0x1593('‮a1', 'kdA6')](_0x8e183, _0x16f16f[_0x1593('‮a2', 'OUQB')]); _0x8e183++) {
            if (_0x447b15[_0x1593('‮a3', 'pRJH')](_0x1593('‮a4', 'lqB9'), _0x447b15[_0x1593('‮a5', 'OUQB')])) {
                var _0x83e6a0 = _0x447b15[_0x1593('‮a6', '@Heq')][_0x1593('‫a7', 'jh%K')]('|'),
                    _0x32d03d = 0x0;
                while (!![]) {
                    switch (_0x83e6a0[_0x32d03d++]) {
                        case '0':
                            var _0x4494ce = _0x447b15[_0x1593('‫a8', 't^J%')](_0x16f16f[_0x8e183]['time'], 0x3e8);
                            continue;
                        case '1':
                            _0x447b15[_0x1593('‫a9', 'Zidc')](_0x5355f0, _0x447b15['TObKV']) ? _0x303875 = _0x1400e9 : '';
                            continue;
                        case '2':
                            var _0x5355f0 = time(_0x4494ce)['split']('\x20')[0x0];
                            continue;
                        case '3':
                            var _0x1400e9 = Math[_0x1593('‫aa', 'Zidc')](_0x447b15[_0x1593('‮ab', 'lqB9')](_0x16f16f[_0x8e183][_0x447b15['nEdGR']], 0x64));
                            continue;
                        case '4':
                            if (_0x447b15['EuCVO'](_0x447b15[_0x1593('‫ac', 'q$H%')](dayDiff, _0x5355f0), 0x1f) && _0x1400e9 <= _0xcb265f[_0x1593('‮ad', '6Qwh')]) {
                                _0xcb265f['price'] = _0x1400e9;
                                _0xcb265f[_0x1593('‮ae', '%E57')] = _0x5355f0;
                            }
                            continue;
                        case '5':
                            _0x447b15[_0x1593('‮af', ']!fw')](_0x5355f0, _0x447b15['yNfJm']) ? _0x199ec1 = _0x1400e9 : '';
                            continue;
                    }
                    break;
                }
            } else { if (_0x50508e) { _0x2ae979 += _0x447b15[_0x1593('‫b0', 't^J%')](_0x447b15[_0x1593('‫b1', '4]QC')], _0x59f556[_0x1593('‮b2', 'qdeo')]); } else { _0x2ae979 += _0x447b15[_0x1593('‮b3', 't^J%')](_0x1593('‮b4', 'D[0X'), _0x59f556[_0x1593('‮b2', 'qdeo')]); } }
        }
        if (_0x447b15['sIMKk'](_0x2f79a1['min'], 0x5f5e0ff)) _0x2f79a1[_0x1593('‮b5', '663P')] = '-';
        if (_0x303875 === 0x0) _0x303875 = '-';
        if (_0x447b15['sIMKk'](_0x199ec1, 0x0)) _0x199ec1 = '-';
        var _0x5d0579 = [];
        _0x5d0579['push']({ 'name': _0x447b15[_0x1593('‮b6', 'kdA6')], 'price': _0x2f79a1[_0x1593('‮b7', 'MPXv')], 'date': _0x2f79a1[_0x1593('‮b8', 'AQJA')], 'diff': _0x447b15[_0x1593('‮b9', 'T7Zb')](priceDiff, _0x2f79a1['max']) });
        _0x5d0579[_0x1593('‫ba', 'pRJH')]({ 'name': _0x447b15['jbaOb'], 'price': _0x2f79a1[_0x1593('‮bb', '$^tX')], 'date': _0x2f79a1['mint'], 'diff': _0x447b15[_0x1593('‫bc', '$^tX')](priceDiff, _0x2f79a1['min']) });
        _0x5d0579['push']({ 'name': '六一八', 'price': _0x303875, 'date': _0x447b15[_0x1593('‫bd', 'kdA6')], 'diff': _0x447b15['AbOpL'](priceDiff, _0x303875) });
        _0x5d0579[_0x1593('‫be', 'kx##')]({ 'name': _0x447b15['owBUG'], 'price': _0x199ec1, 'date': _0x1593('‫bf', '2bIO'), 'diff': _0x447b15['WsIxb'](priceDiff, _0x199ec1) });
        var _0x5cb665 = '';
        for (var _0x8e183 = 0x0; _0x447b15[_0x1593('‮c0', 'MPXv')](_0x8e183, _0x5d0579['length']); _0x8e183++) { if (_0x447b15[_0x1593('‮c1', 'Zidc')] === _0x1593('‫c2', 'Z[yY')) { if (_0x447b15[_0x1593('‫c3', 'D[0X')](typeof old, 'number')) return ''; var _0x2303e1 = _0x447b15[_0x1593('‮c4', '@Heq')](old, now); if (Math['abs'](_0x2303e1) < 0.1) { _0x2303e1 = 0x0; } if (_0x447b15[_0x1593('‫c5', 'EJt%')](_0x2303e1, 0x0)) return ''; return _0x447b15[_0x1593('‫c6', 'WUnd')](_0x2303e1, 0x0) ? _0x447b15['qqqoa']('↓', Math[_0x1593('‮c7', 'AQJA')](_0x2303e1)) : _0x447b15['rnQxW']('↑', Math[_0x1593('‮c8', ']!fw')](Math[_0x1593('‮c9', '6Qwh')](_0x2303e1))); } else { if (!(_0x447b15[_0x1593('‮ca', 'lqB9')](_0x5d0579[_0x8e183][_0x1593('‮cb', 't^J%')], '-') || _0x5d0579[_0x8e183]['price'] === 0x0)) _0x5cb665 += _0x447b15[_0x1593('‫cc', 't^J%')](_0x447b15[_0x1593('‫cd', 'Lvyu')](_0x5d0579[_0x8e183][_0x1593('‫ce', '$Ps&')] + '：', _0x5d0579[_0x8e183][_0x447b15['nEdGR']]), '\x0a'); } }
        _0x2ae979 += _0x5cb665;
        _0x2ae979 += _0x447b15[_0x1593('‫cf', 'ZiSQ')];
    } else if (_0x5e4a3f && _0x5e4a3f[_0x1593('‮d0', 'Lvyu')]) {
        for (var _0x8e183 = 0x0; _0x447b15['CPFGH'](_0x8e183, _0x5e4a3f[_0x1593('‮d1', 'EJt%')][_0x1593('‮d2', '%E57')]); _0x8e183++) { var _0x28036b = _0x5e4a3f[_0x1593('‫d3', 'qdeo')][_0x8e183]; if (_0x447b15[_0x1593('‮d4', 'ZSjm')](_0x28036b[_0x1593('‫d5', '1sX^')][_0x1593('‮d6', 'pRJH')](_0x447b15[_0x1593('‮d7', 'qdeo')]), -0x1)) { if (_0x447b15[_0x1593('‫d8', 'EJt%')](_0x1593('‫d9', 'H]Ev'), _0x447b15[_0x1593('‮da', '@JFW')])) { _0x2ae979 += _0x447b15['rnQxW'](_0x447b15[_0x1593('‮db', '4]QC')], _0x59f556[_0x1593('‫dc', '%E57')]); } else { _0x2ae979 += _0x447b15[_0x1593('‮dd', '663P')](_0x447b15['qVTdk'](_0x447b15[_0x1593('‫de', 'i*ZX')], _0x28036b[_0x1593('‫df', 'Z[yY')]), '\x0a'); } } else if (_0x447b15['ChmqT'](_0x28036b[_0x1593('‫e0', 'J66I')][_0x1593('‫e1', 'D[0X')]('11'), -0x1)) { if (_0x447b15[_0x1593('‮e2', ']]PH')]('TVzJl', 'uidaw')) { _0x2ae979 += _0x447b15[_0x1593('‫e3', '1sX^')](_0x447b15['zSSIl'], _0x59f556['short']); } else { _0x2ae979 += _0x447b15[_0x1593('‮e4', ']]PH')](_0x447b15[_0x1593('‫e5', '7$wv')] + _0x28036b[_0x1593('‮e6', 'qdeo')], '\x0a'); } } }
        _0x2ae979 += _0x447b15[_0x1593('‫e7', 'q$H%')];
    }
    _0x2ae979 += _0x447b15[_0x1593('‫e8', 'qdeo')](_0x1593('‫e9', ']]PH'), now) + '\x0a';
    if (_0x5e4a3f) {
        if (_0x447b15[_0x1593('‮ea', 'Lvyu')](_0x1593('‫eb', 'aE9n'), _0x1593('‫ec', 'pRJH'))) {
            var _0x5296cd = _0x447b15[_0x1593('‮ed', 'f(aS')]['split']('|'),
                _0x1bc6a2 = 0x0;
            while (!![]) {
                switch (_0x5296cd[_0x1bc6a2++]) {
                    case '0':
                        if (_0x447b15[_0x1593('‮ee', 'aE9n')](_0x447b15['dLcWb'](_0x398b1f, 0x17), 0x0) && _0x398b1f) { _0x31dab0 = !![]; }
                        continue;
                    case '1':
                        if (_0x447b15['dLcWb'](_0x51da9c, 0x7) == 0x0 || force) { _0x398b1f = _0x447b15[_0x1593('‫ef', 'ekB8')](_0x3fcd44); }
                        continue;
                    case '2':
                        var _0x398b1f = get('random');
                        continue;
                    case '3':
                        if (!_0x447b15['mDGED'](get, 'jd_spy_home')) { return ![]; }
                        continue;
                    case '4':
                        var _0x31dab0 = ![];
                        continue;
                    case '5':
                        return _0x31dab0;
                    case '6':
                        _0x398b1f = _0x447b15['mDGED'](parseInt, _0x398b1f);
                        continue;
                    case '7':
                        var _0x51da9c = new Date()[_0x1593('‮f0', '$^tX')]();
                        continue;
                }
                break;
            }
        } else { _0x2ae979 += _0x447b15[_0x1593('‫f1', 'q$H%')](_0x447b15[_0x1593('‮f2', 'Y^Ra')](_0x447b15[_0x1593('‫f3', 'WUnd')]('(', _0x5e4a3f[_0x1593('‫f4', 't^J%')]), ')'), '\x0a'); }
    }
    if (_0x59f556 && _0x59f556[_0x1593('‮f5', '@Heq')] != _0x59f556[_0x1593('‮f6', 'WUnd')] && _0x59f556[_0x1593('‫f7', 't^J%')]) _0x2ae979 += _0x447b15[_0x1593('‮f8', 'kx##')] + _0x59f556[_0x1593('‫f9', '2bIO')] + '\x0a';
    if (_0x59f556 && _0x59f556[_0x1593('‫fa', '2bIO')]) {
        _0x2ae979 += '\x0a';
        if (_0x48cd99) {
            if (_0x447b15['JNwJh'](_0x447b15['IWnUy'], _0x447b15[_0x1593('‮fb', 'qdeo')])) {
                if (_0x2cf814) { if (_0x50508e) { if ('VOqLX' !== _0x447b15['nluzN']) { _0x447b15[_0x1593('‫fc', 'OUQB')](set, _0x447b15[_0x1593('‫fd', '@JFW')], r); return +r; } else { _0x2ae979 += _0x447b15['HqGbO'](_0x447b15['PZoow'], _0x59f556['short']); } } else { if (_0x447b15['cWxEB'](_0x447b15['YVRFA'], _0x447b15[_0x1593('‮fe', 'R#xJ')])) { _0x2ae979 += _0x447b15[_0x1593('‮ff', 'q$H%')](_0x447b15[_0x1593('‮100', '$Ps&')], _0x59f556[_0x1593('‮101', '!yu5')]); } else { if (_0x50508e) { _0x2ae979 += _0x447b15['QfFBi'](_0x447b15[_0x1593('‫102', 'O&K2')], _0x59f556[_0x1593('‫103', 'gWcR')]); } else { _0x2ae979 += _0x447b15['pAJTa'](_0x447b15[_0x1593('‫104', 'Zidc')], _0x59f556['short']); } } } } else {
                    if (_0x50508e) {
                        if (_0x1593('‮105', 'kdA6') !== 'lAdfZ') {
                            for (var _0x52fcf6 = 0x0; _0x52fcf6 < _0x5e4a3f[_0x1593('‮106', '2bIO')][_0x1593('‮107', '7$wv')]; _0x52fcf6++) { var _0x4d4382 = _0x5e4a3f['promo_days'][_0x52fcf6]; if (_0x4d4382['show']['indexOf'](_0x1593('‫108', 'i*ZX')) != -0x1) { _0x2ae979 += _0x447b15[_0x1593('‮109', 'i*ZX')](_0x447b15[_0x1593('‮10a', 'lqB9')] + _0x4d4382[_0x1593('‮10b', 'T7Zb')], '\x0a'); } else if (_0x447b15['dzDET'](_0x4d4382['show'][_0x1593('‮10c', '6Qwh')]('11'), -0x1)) { _0x2ae979 += _0x447b15['gUtTl'](_0x447b15[_0x1593('‫10d', 'Zidc')], _0x4d4382['price']) + '\x0a'; } }
                            _0x2ae979 += _0x1593('‮10e', 'EJt%');
                        } else { _0x2ae979 += _0x447b15[_0x1593('‫10f', 'OUQB')](_0x447b15[_0x1593('‮110', '$^tX')], _0x59f556[_0x1593('‮111', '1sX^')]); }
                    } else { _0x2ae979 += _0x447b15['wroKh'](_0x447b15[_0x1593('‫112', ']!fw')], _0x59f556['short']); }
                }
            } else { diff = 0x0; }
        } else { if (_0x50508e) { _0x2ae979 += _0x447b15['LmSnZ'](_0x447b15[_0x1593('‮113', 'J66I')], _0x59f556['short']); } else { _0x2ae979 += _0x447b15[_0x1593('‫114', 'S14v')](_0x447b15[_0x1593('‫115', '@Heq')], _0x59f556['short']); } }
    }
    _0x447b15[_0x1593('‮116', 'LlRA')](sendText, _0x2ae979[_0x1593('‮117', 'Lvyu')]('\x0a'));
}

function time(_0x3ad131) { var _0x35b865 = { 'LkQEZ': function(_0x903c4, _0x3d6c07) { return _0x903c4 + _0x3d6c07; }, 'LuUCs': function(_0x2b869a, _0x5adc56) { return _0x2b869a * _0x5adc56; } }; if (!_0x3ad131) {+new Date(); } var _0x359701 = new Date(_0x35b865['LkQEZ'](_0x3ad131, _0x35b865[_0x1593('‮118', 'O&K2')](0x8 * 0xe10, 0x3e8))); return _0x359701[_0x1593('‫119', 'WUnd')]()[_0x1593('‮11a', 'AQJA')](0x0, 0x13)[_0x1593('‫11b', 'Zz^R')]('T', '\x20')[_0x1593('‫11c', 'f(aS')]('\x20')[0x0][_0x1593('‮11d', 'ZiSQ')](/\./g, '-'); }

function dayDiff(_0x42a2af) { var _0x5c23d4 = { 'hYOHg': function(_0x5e4218, _0xb7326c) { return _0x5e4218 * _0xb7326c; }, 'fmmVy': function(_0x3164bb, _0x5e0a2e) { return _0x3164bb * _0x5e0a2e; } }; return parseInt((new Date() - new Date(_0x42a2af)) / _0x5c23d4[_0x1593('‫11e', 'pRJH')](_0x5c23d4['fmmVy'](0x3e8 * 0x3c, 0x3c), 0x18) + ''); }

function priceDiff(_0x201191) { var _0x2cb9c = { 'hVftq': function(_0x4b8ff7, _0x196fc4) { return _0x4b8ff7 !== _0x196fc4; }, 'PzfEY': _0x1593('‮11f', 'G&AV'), 'BboUd': function(_0x44a352, _0x38355a) { return _0x44a352 - _0x38355a; }, 'Febkk': function(_0x29ea57, _0x382921) { return _0x29ea57 < _0x382921; }, 'TPEmy': _0x1593('‫120', 'qdeo'), 'BnTQD': function(_0x1e92e9, _0x780476) { return _0x1e92e9 === _0x780476; }, 'BzRnF': function(_0x38afee, _0x358812) { return _0x38afee > _0x358812; }, 'qzpOc': function(_0x45d402, _0x1a345f) { return _0x45d402 + _0x1a345f; } }; if (_0x2cb9c['hVftq'](typeof _0x201191, _0x2cb9c[_0x1593('‫121', 'aE9n')])) return ''; var _0x4c856d = _0x2cb9c[_0x1593('‫122', 'f(aS')](_0x201191, now); if (_0x2cb9c['Febkk'](Math[_0x1593('‫123', '$Ps&')](_0x4c856d), 0.1)) { if (_0x2cb9c['TPEmy'] !== _0x2cb9c[_0x1593('‮124', '$Ps&')]) { auth_fanli = !![]; } else { _0x4c856d = 0x0; } } if (_0x2cb9c[_0x1593('‮125', 'R#xJ')](_0x4c856d, 0x0)) return ''; return _0x2cb9c[_0x1593('‮126', 'WUnd')](_0x4c856d, 0x0) ? '↓' + Math[_0x1593('‫127', 'q$H%')](_0x4c856d) : _0x2cb9c[_0x1593('‫128', 'ZiSQ')]('↑', Math[_0x1593('‫129', ']]PH')](Math[_0x1593('‮12a', '@Heq')](_0x4c856d))); }

function space(_0x218679, _0x428a58) { var _0x3f3471 = { 'fBWFO': _0x1593('‮12b', '663P'), 'EXCGC': function(_0x31cd13, _0x500074) { return _0x31cd13 < _0x500074; }, 'OYgQz': function(_0x9cbb93, _0x4d90e4) { return _0x9cbb93 - _0x4d90e4; }, 'BMlcY': function(_0xddc024, _0xba9a85) { return _0xddc024 !== _0xba9a85; }, 'afbIL': 'LlkPe' }; var _0x162f83 = ''; for (var _0x1465e8 = 0x0; _0x3f3471[_0x1593('‮12c', 'OUQB')](_0x1465e8, _0x3f3471[_0x1593('‫12d', '!yu5')](_0x428a58, (_0x218679 + '')[_0x1593('‫7a', 'aE9n')])); _0x1465e8++) { if (_0x3f3471['BMlcY'](_0x3f3471[_0x1593('‫12e', 'kdA6')], _0x1593('‫12f', ']]PH'))) { _0x162f83 += '\x20'; } else { rt += _0x3f3471['fBWFO'] + info['short']; } } return _0x162f83; }
main();;
_0xod7 = 'jsjiami.com.v6';