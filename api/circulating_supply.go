package api

import (
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

func (a *API) updateCirculatingSupply() {
	c := a.getCirculatingSupply()
	a.CirculatingSupply = c

	t := time.NewTicker(10 * time.Minute)

	for {
		select {
		case <-t.C:
			c := a.getCirculatingSupply()

			a.CirculatingSupply = c
		}
	}
}

func (a *API) getCirculatingSupply() decimal.Decimal {
	return totalSupply.Sub(a.getBalances())
}

func (a *API) getAddressBalance(address string, name string) decimal.Decimal {
	var log = log.WithField("target", name)

	log.Infof("fetching value locked")
	val, err := a.xyz.BalanceOf(nil, common.HexToAddress(address))
	if err != nil {
		log.Error(errors.Wrap(err, "could not fetch value locked in community vault"))
	}

	valDec := decimal.NewFromBigInt(val, -18)

	log.WithField("value", valDec).Info("done")

	return valDec
}

func (a *API) getBalances() decimal.Decimal {
	var addrs = []string{
		// vesting contracts
		"0x5520A5EF5D5093b8Fd650a7a9602ed4418CA2060",
		"0x965E039e0B844d7Dd441Ce25B08d0073738257F0",
		"0x3a50B410aF330BDdDEdB497f33CA2C81FBB74d8F",
		"0x0FAea946a18F4FeEAd8F33833c0CE3cCf40bC70F",
		"0x709d44ADaddA2004C0b0681421d59C321DF6CEf8",
		"0xDD76eA5A131982fC99480A46f928F8a5F368de6f",
		"0x386153eBB6Ce46199fb6D675100ceAB2172E47ab",
		"0x2e7F9B3e12510D7a0F41a8Bdb3975d27569Ea347",
		"0x2F50076eBbBE0430957bE1CEf0DB487f2BB65676",
		"0xbF1AC089bbAC2c749a7c011BEBCBbe0628D1E194",
		"0x14dbFf41743D389a69fD208a7F41021Cc8FB9ccf",
		"0x1568b2eFd0e98bdBABCA879F2a778f7EaA2323dE",
		"0xDc6A0C905083e1a9417468180790150aF914e2C8",
		"0xF38Cc32A49435E5c272362193580177b31B13aBF",
		"0x1133B8Ed56629b44d29eC12D0E323A30d07f6490",
		"0x846f0Da84167a161cD8132b7934D765Ed3d54451",
		"0x191E56Ab8FEbf1CAF14859C7720B201eEcBD5a3e",
		"0x3395c2E42599CD88342A9802115C33E2D8ec63E0",
		"0xc0e835b13fB4Fa5d524B2C02Fe74e66d21CA275A",
		"0xF8A3D840A522A281C432F3d4c162a15569f97454",
		"0x454E0eDceE4a1D71079F608C834753aDf7E7bDA9",
		"0xB19683cC2a11cCFF68cbB9c0D25F99b9B3437baE",
		"0x4D116Df3420b0281386F7157dB7303c1d85D1B78",
		"0x0ed3B3E9CEC482A42f77A3cB5e494892D2b9a8C3",
		"0x8957ecC1820d13F9169B56Ab933a47fDe00ce2b8",
		"0x640D370f813c444a86D0be7216E8374993E0FF7a",
		"0x10f27E99b6b0c30eE19035C4B5AeC556178c35Cc",
		"0xf4a302C5e6A0d975cC4C613d4A9DC0D754EC56b9",
		"0xDd6506A5619A0a0A2E779FFeB237f808c3B1D4d7",
		"0x559896623DDE78808Bbfc18d83003805A44D0ec2",
		"0xb9C6A671beaB22Ca56EF6E92c67795E9Ce7D685C",
		"0xDb87d27d47Fe0Bb4a96C72a8dd02B02E4cE735C0",
		"0x72Af11a5fBE278a225209D7aC055f51A2f8Fa1B4",
		"0x04f1dcc48dc36958fF441893c3d30F0aB5e521F2",
		"0x853448412a14cb3B324BFc5a2Aeb120E0b7AC342",
		"0x1531d969050a3851b3C0B82006ff1259D18a988E",
		"0x6D85379Da4D0ff46677DBa31B602AF66D4E7123F",
		"0x642BbE39FE70450267D206d5151998c49C6E050E",
		"0x09840046a34cAaCC193aFceA76b81426F67Fb963",
		"0x7Ba69955B1565aBBDaE1D0E1d75C567115E9290C",
		"0x0c19Ef41fA36433D380b09aC6Cc44432eA268D2C",
		"0xf386CcBC27A5c24694ECa0a7e98E0D9BEBaEf64A",
		"0x16d85C49889f1E404d53f1c7F4c182BBFAf7674D",
		"0xf8e6c795c41592548d2B955451ba5C85a669024f",
		"0xC77719515E7304fDBb25d6bc4A0E58dAe4955961",
		"0x8a77Cd92851c4d52d79302BB46BA131d89314007",
		"0x15c93650f483247d6380a404C0d0879f2344AE0a",
		"0x2f25C243bC59f4f4DfC52a0828BC7EEFdcdADFE9",
		"0x90ea42354892aebF1676Cfd5DF823E8951C68aCa",
		"0x0D90a3a4275bE9dC5cb95D0F83e3e24c0669601C",
		"0x230407FE2F5cb31d5aC6dB1F1F8b8B5b275C4f20",
		"0x55a8e686dE4530378636CeF5F892678Da1Bc2976",
		"0x7d0fe8d45EdcE3D95e410f224E958dB83590EE4A",
		"0x4808A637016a1a25a3602EAF40fD1C6D6aDb7486",
		"0x3BC624Ed8c8D3B620B46939092C238583A938a76",
		"0xbb4b3FAf0c70eD12BE6752446d88548f91E06173",
		"0x7bD02B53C3e8d400379176f310EfF63FcAB9aC45",
		"0x85409fFE72353EB685f6F575A8615Cfc9BD019F3",
		"0x0cDA7Fecc1e1c32426833fa7eca097Ad4906157A",
		"0x852077C6D03ea6c3848F708C5FE1f06156D69B9c",
		"0x15B789A6eF80836c10233A1153740cdB95993386",
		"0x64dDFc6Cb4D81e06EEf5271fABfAC9Bb103D6C3c",
		"0x2181d9E6508259cA7B0d9C3CFB445305a8D1D1BE",
		"0x5384b6AE4Ec81822c07F182ca9E1282520450630",
		"0x73e68BFD4c173021De72A312A5E3140e3D34930B",
		"0xe9d10f1559dE30dcD48FB22AB1151a40b0f9B433",
		"0x6cF3Da61bd2D9482bA74e885f752B32301AAa110",
		"0x29C9563bd9A79AdA30A4C232f4ED45dA9899f9AE",
		"0xCE49E21C6107589068229D10d30ae8750362b7B4",
		"0xDf206fBAAb393ED9053653aBa0d557C769dBA130",
		"0x3aC7c6E8753fC52DD6C1022a555cb02c4e94cE29",
		"0x8c9596046b81E8747a9c2cBbE4E916a2e1690efc",
		"0x7058045aB6CD975eC81c47ceb4AbFFd04a0A68F3",
		"0xE24ddE54A55585Dd456c0BdF64C4C7DD5e0DD847",
		"0x854a238C133feD2A40a875729C679A1341E9e9CA",
		"0x3d282D2861FeFDb0da6c038B05435beDA49a9b86",
		"0xf081Eec13bE8D99eAFCc2e6cED6091ce7A18B764",
		"0xDD510bc858Eb0cfBA39e3e2f328F6eF91fB50d3E",
		"0xF25Ecb13715E9D0BFd9db4E6A499B0ff0CF05b3e",
		"0x45b90eD252646Dbf744a1f6aB315cf19239C6039",
		"0xEd15b185fB42003B03b89fa185ee6a6de0778618",
		"0x26e1164582258851889A0256D4BBC11c47DEEAD1",
		"0xd292c518636E37Ad622a1F96b1B55535CC97E95C",
		"0x11DD9d239D9b047D3b8EC870f40ad207c1dcb414",
		"0xeF243965aC864423b43335d0a4C23067755e7F4C",
		"0xA9a534eA789B0FA20266FB6a036f5Aa318B61e85",
		"0x64f9F22Fa7f761C21d169Df2b3343B18a9b50a15",
		"0xfFa1cd61e21B16B78E3Fd3481267e67d3316980E",
		"0x3363FAC31fdBd00707C339Bb117992B0fFf8a5dC",
		"0xd5ca084B66b4C6279E75E0A79246140B189C6e32",
		"0xa2E754Ab02cd8b9780E872C3dd1B7735373351e5",
		"0xa3e080D9ace354C3A6C0E3B2B0a9f23298d94353",
		"0x783CF9Cf627c87aC5790FC6A9d78E348F890E42F",
		"0xE83Aab72F4e016Cde12aDA331d28a4DaFa3B74f3",
		"0x1cC6316eF595f7E299065ea33c74D11a1Cf1eEa3",
		"0x3e26B62df967DA08eba64B9800CC422A8c79caa8",
		"0x1fdc8D2E6b73329677011d9f2a1c2Ec4D9dc00f7",
		"0x9D063Ca9C7F92089614DA196b2d88afDf343B0C7",
		"0x3d44A19AE6217381bdb782505b6d2b0642Aa4a13",
		"0xe050e0696fB9EB7C2A424025AAd3BDCC5B877793",
		"0x6C132088dEF7c6b6c2f245362C3c14f9A7eBef85",
		"0xD6173B942766E250CC87e1841D572f8C214C3CC4",
		"0xA1bbaDB01f9B85176b0CeD3e344B7a8141043841",
		"0x7336aFA7198e913376d3A34563EC034BA6F2F191",
		"0x40FCB910B9FefC9c6402750b733E469B4BCD44D1",
		"0x6592FBeEB3de5aCa466167e49eC66E0258a4115b",
		"0x4955Ca593e557674B4a9FeE2466bE9295B8242dB",
		"0xAD6EB0fb104091FA71E2C3861eCCbc3A6916C75a",
		"0xf4aB1932c85B8FD52C7b605E108E1e01A069dA18",
		"0x4d0640B22C701Ff95d6F66FaA1Bddd76c1c9E8AB",
		"0x889A70B0F8B49340812FEE9f1ba229225BC47b95",
		"0xdfb5124b60d03cb33405D86D2906A52804955f27",
		"0xC94c554BFCc7D5FdbABD46817D235892E2D62201",
		"0x786Ee8962A30D5bee8A848e98CF75c6a01eb190A",
		"0xD397951A94342B8A9Bd5Dc6518E304E75F225882",
		"0xd756Ea477f5E14E9971faE626d1251Efd07FE2A6",
		"0xb43B25CE9cBd9a6073Fbe5B3Ffc0aEF8f49373AC",
		"0x6B91B1893e9be183bAEDe5aD32c9418ec5fD56EB",
		"0x6b0F2F581dF2A857d562CFCf9D89A0f14D7510f7",
		"0xFeA490228B441F03076A69BA5e084Ae2C6e057E1",
		"0x2B014d5869a827C86a4dC68ec442c38E414F88D1",
		"0x99CD8C52DA9BCe099FC443ee968e00F3bBbC31f5",
		"0x51022e453E2b779FF943b4e95b23e3fD2946D607",
		"0x0321D99EF353FD25A29b0805125F9203C75352f8",
		"0xA5D9B86abD9C5Be4d82aCE958DDe9E6d1EA315dB",
		"0xCA46FD4E7c14c3c095DC0E226ABcb24b511DeD64",
		"0xB7D4079d6890F7fdCA7be2329Cf985FbfD6209a7",
		"0xc387DCE727163beD366DC3C8A681974dF0b4eF87",
		"0x0C638d2ca0f41A01da831f57087cB09a4438531e",
		"0x63ed6e696C5783409c396ddCa3CED1F8e9b7c206",
		"0xdeD5C1693A4D1362C9Ae250d9b4b10380Ac9f690",
		"0xB9adef4F1938884F2B82095c6e837b77e5D39af1",
		"0x14F768075368d21095416dC9504c91F77288CB34",
		"0xD480c9325029528740376cE567e8b39bF27A98fc",
		"0x07AFd4D6B73Fa8DE268868f398D4A50f1A10a7cc",
		"0x7D2691f5Ceb0f423d9D5f0eE411B29AE56f8AeDa",
		"0xAc545197817c592B4a3202C561e6f6EE27a3A36c",
		"0xe7CA02be246E0BCb649c8ab2A0Ff42696DD2b0f9",
		"0x646D08Da23adB3D71004C44Ae99caa2EBfB43a0c",
		"0x6b96E1df9529Ec1A39514767B15167b2aD1bB61F",
		"0x19B1351891e8A5589eF5Df4af6Cb0a0f3CEE853E",
		"0x910Fa90Bb0335f7d5539775fd9A494e64E1A0785",
		"0x4cc488dA00DD587F1d664c7AF5345A9D6cFe1807",
		"0x4243b297cCfbfC6F7615fb4B4e6E3373ECF2955b",
		"0x25ea79601432559245d938372735e341dd10B242",
		"0xeF2abBD50D83E40219af4f66c50e6789A70C26e5",
		"0x9C9cBC1B82672AAbb26D76b1E09035a354b9991a",
		"0x036654cf1F0FEa99d58ef045e5BaD6587EFDeD99",
		"0xe6fA7Cc2AB8cAb5f77A752D8161A0E9978ea2935",
		"0xBb773d84Bc0e24C66AbF9C94dBdCA278dE0Ab3c6",
		"0x82bcD4e405FdC2851b400C5fB36bbddF66a9FF83",
		"0x222C142D8ac0aDF27d20801464DE93431c00f822",
		"0xd0bD07a31411423D10752a03C18F0fD26Eb6937d",
		"0x18c5203287Bdd3A7098d11Ad7f67CB1E66E7e7c5",
		"0x074BBf1d2169370201D514B912d575cD5f461249",
		"0xc25045167526951c1Dc50034EfC6Ea71E2D60154",
		"0xd0F6A858728E00C841c978Ade75E288be24C85F6",
		"0x3F2De6BadBA37f387d6A02Be77671e1314aCcE7c",
		"0x7800Cde240Aeb966be7B8D77F41d8ebc1431751E",
		"0xAAFF4C18e871E26EdaA5f2443382d4f3f1917c5B",
		"0x11131648a7961439F5d92AF074d03EADFB5Ee755",
		"0xeA0c0Cc51608207AAc3975C948066c6cef41D036",
		"0xEaD154C123C345923265B29F7945094575CEd8bA",
		"0x8315124C30CcFd55a204A34bAdD7e20FECb45765",
		"0x5F6A1e45Df5eafFA82148c23115ba989BC2749Cc",
		"0xDf42c0E8C1d80f3B02F298A8Ab00B9B68ccDd5B3",
		"0xD847d728964F0bfE5246F487cC56701A4d5608C5",
		"0x54AE049C62935d6D7b40b7135717E1467a2B49cB",
		"0x3736cEB4A896f7DddDE7df29E3958DC74EF10214",
		"0x8B53814a884309Dc7e280907582ecABb5a8CCcb4",
		"0x41092c27ca1cE17Ed8C4Caa34c4469CDbC75D757",
		"0xF2bDbeF1240Ba2052E81546197F662907bc27c68",
		"0xf97664376416E9379f2354DB444BFE3f00B6936b",
	}

	var total decimal.Decimal
	for _, addr := range addrs {
		total = total.Add(a.getAddressBalance(addr, addr))
	}

	return total
}