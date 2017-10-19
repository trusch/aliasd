package storage_test

import (
	"os"

	. "github.com/trusch/aliasd/storage"
	"github.com/trusch/storage/engines/meta"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Storage", func() {
	BeforeEach(func() {
		os.MkdirAll("/tmp/testdb", 0755)
	})
	AfterEach(func() {
		os.RemoveAll("/tmp/testdb")
	})

	It("should be possible to set/get values", func() {
		db, err := meta.NewStorage("leveldb:///tmp/testdb")
		Expect(err).NotTo(HaveOccurred())
		store, err := NewStorage(db, "default")
		Expect(err).NotTo(HaveOccurred())
		Expect(store.Set("foo", "bar")).To(Succeed())
		val, err := store.Get("foo")
		Expect(err).NotTo(HaveOccurred())
		Expect(val).To(Equal("bar"))
		Expect(store.Close()).To(Succeed())
	})

	It("should be possible to set/get values after restart", func() {
		db, err := meta.NewStorage("leveldb:///tmp/testdb")
		Expect(err).NotTo(HaveOccurred())
		store, err := NewStorage(db, "default")
		Expect(err).NotTo(HaveOccurred())
		Expect(store.Set("foo", "bar")).To(Succeed())
		Expect(store.Close()).To(Succeed())

		db, err = meta.NewStorage("leveldb:///tmp/testdb")
		Expect(err).NotTo(HaveOccurred())
		store, err = NewStorage(db, "default")
		Expect(err).NotTo(HaveOccurred())
		val, err := store.Get("foo")
		Expect(err).NotTo(HaveOccurred())
		Expect(val).To(Equal("bar"))
	})
})
