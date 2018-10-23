package oracle_test

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/tokencard/ethertest"
)

var _ = Describe("removeTokens", func() {

	Context("When tokens are supported", func() {
		BeforeEach(func() {

			tokens := []common.Address{common.HexToAddress("0x1"), common.HexToAddress("0x2"), common.HexToAddress("0x3")}
			tx, err := oracle.AddTokens(
				controller.TransactOpts(),
				tokens,
				stringsToByte32(
					"OMG",
					"EOS",
					"TKN",
				),
				[]*big.Int{
					calculateMagnitude(big.NewInt(18)),
					calculateMagnitude(big.NewInt(18)),
					calculateMagnitude(big.NewInt(8)),
				},
				big.NewInt(20180913153211),
			)
			Expect(err).ToNot(HaveOccurred())
			be.Commit()
			Expect(isSuccessful(tx)).To(BeTrue())
		})
		Context("When called by the controller", func() {

			Context("When removing a supported token", func() {

				var tx *types.Transaction
				BeforeEach(func() {
					var err error
					tx, err = oracle.RemoveTokens(controller.TransactOpts(), []common.Address{common.HexToAddress("0x2")})
					Expect(err).ToNot(HaveOccurred())
					be.Commit()
				})

				It("Should succeed", func() {
					Expect(isSuccessful(tx)).To(BeTrue())
				})

				It("Should emit a TokenRemoval event", func() {
					it, err := oracle.FilterRemovedToken(nil)
					Expect(err).ToNot(HaveOccurred())
					Expect(it.Next()).To(BeTrue())
					evt := it.Event
					Expect(it.Next()).To(BeFalse())
					Expect(evt.Token).To(Equal(common.HexToAddress("0x2")))
				})

				It("Should update the tokens map", func() {
					token, err := oracle.Tokens(nil, common.HexToAddress("0x2"))
					Expect(err).ToNot(HaveOccurred())
					Expect(token.Exists).NotTo(BeTrue())
					Expect(token.LastUpdate.String()).To(Equal("0"))
					Expect(token.Symbol).To(Equal(""))
					Expect(token.Magnitude.String()).To(Equal("0"))

					//the other tokens should remain unchanged
					token, err = oracle.Tokens(nil, common.HexToAddress("0x1"))
					Expect(err).ToNot(HaveOccurred())
					Expect(token.Exists).To(BeTrue())
					Expect(token.LastUpdate.String()).NotTo(Equal("0"))
					Expect(token.Symbol).To(Equal("OMG"))
					Expect(token.Magnitude.String()).To(Equal(calculateMagnitude(big.NewInt(18)).String()))
				})
			})

			Context("When removing all supported tokens", func() {

				var tx *types.Transaction
				BeforeEach(func() {
					var err error
					tokens := []common.Address{common.HexToAddress("0x1"), common.HexToAddress("0x2"), common.HexToAddress("0x3")}
					tx, err = oracle.RemoveTokens(controller.TransactOpts(), tokens)
					Expect(err).ToNot(HaveOccurred())
					be.Commit()
				})

				It("Should succeed", func() {
					Expect(isSuccessful(tx)).To(BeTrue())
				})

				It("Should emit 3 TokenRemoval events", func() {
					it, err := oracle.FilterRemovedToken(nil)
					Expect(err).ToNot(HaveOccurred())
					Expect(it.Next()).To(BeTrue())
					evt := it.Event
					Expect(it.Next()).To(BeTrue())
					Expect(evt.Token).To(Equal(common.HexToAddress("0x1")))
					evt = it.Event
					Expect(it.Next()).To(BeTrue())
					Expect(evt.Token).To(Equal(common.HexToAddress("0x2")))
					evt = it.Event
					Expect(it.Next()).To(BeFalse())
					Expect(evt.Token).To(Equal(common.HexToAddress("0x3")))
				})
			})

			Context("When removing all supported tokens but a duplicate is passed", func() {
				It("Should fail", func() {
					tokens := []common.Address{common.HexToAddress("0x3"), common.HexToAddress("0x2"), common.HexToAddress("0x1"), common.HexToAddress("0x2")}
					tx, err := oracle.RemoveTokens(controller.TransactOpts(WithGasLimit(300000)), tokens)
					Expect(err).ToNot(HaveOccurred())
					be.Commit()
					Expect(isGasExhausted(tx, 300000)).To(BeFalse())
					Expect(isSuccessful(tx)).To(BeFalse())
				})
			})

			Context("When removing at least one unsupported token", func() {
				It("Should fail", func() {
					tokens := []common.Address{common.HexToAddress("0x1"), common.HexToAddress("0x0"), common.HexToAddress("0x2")}
					tx, err := oracle.RemoveTokens(controller.TransactOpts(WithGasLimit(300000)), tokens)
					Expect(err).ToNot(HaveOccurred())
					be.Commit()
					Expect(isGasExhausted(tx, 300000)).To(BeFalse())
					Expect(isSuccessful(tx)).To(BeFalse())
				})
			})

			Context("When trying to remove #0 tokens (empty list)", func() {

				var tx *types.Transaction
				BeforeEach(func() {
					var err error
					tx, err = oracle.RemoveTokens(controller.TransactOpts(), nil)
					Expect(err).ToNot(HaveOccurred())
					be.Commit()
				})

				It("Should succeed", func() {
					Expect(isSuccessful(tx)).To(BeTrue())
				})

				It("Should emit no Removed Token event", func() {
					it, err := oracle.FilterRemovedToken(nil)
					Expect(err).ToNot(HaveOccurred())
					Expect(it.Next()).To(BeFalse())
				})
			})

		})

		Context("When called by a random address", func() {
			It("Should fail", func() {
				tx, err := oracle.RemoveTokens(randomAccount.TransactOpts(WithGasLimit(300000)), []common.Address{common.HexToAddress("0x1")})
				Expect(err).ToNot(HaveOccurred())
				be.Commit()
				Expect(isGasExhausted(tx, 300000)).To(BeFalse())
				Expect(isSuccessful(tx)).To(BeFalse())
			})
		})
	}) //when tokens are supported

})