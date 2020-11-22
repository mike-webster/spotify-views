package data

import (
	"context"
	"fmt"
	"testing"

	"github.com/bmizerany/assert"
)

func TestEncryption(t *testing.T) {
	key := fmt.Sprint("!;$W=T3rYHXpB'K^")
	ctx := context.WithValue(context.Background(), ContextSecurityKey, key)
	val := "AQAe2FljEaseIixwMLbADRq_2xtOQwNNs5K_MmwTYRq32IzwrJ7aqlqRS9DuRsGQoD8ZxUIeKaHL4Er1fDkL7HtY1kzeUckE3FMAZmnF6zBpjf91IjbXYcOorDOnoE-_w5c"

	encrypted, err := encrypt(ctx, val)
	t.Run("EncryptSuccess", func(t *testing.T) {
		assert.Equal(t, err, nil)
		assert.NotEqual(t, encrypted, val)
	})

	dec, err := decrypt(ctx, encrypted)
	t.Run("DecryptSuccess", func(t *testing.T) {
		assert.Equal(t, nil, err)
		assert.Equal(t, dec, val)
	})

	t.Run("CheckingSubsequebntDecrypt", func(t *testing.T) {
		e3, err := decrypt(ctx, encrypted)
		assert.Equal(t, nil, err)
		assert.Equal(t, val, e3)
	})

	t.Run("CheckingSubsequebntDecrypt", func(t *testing.T) {
		e3, err := decrypt(ctx, encrypted)
		assert.Equal(t, nil, err)
		assert.Equal(t, val, e3)
	})

	t.Run("CheckingSubsequebntDecrypt", func(t *testing.T) {
		e3, err := decrypt(ctx, encrypted)
		assert.Equal(t, nil, err)
		assert.Equal(t, val, e3)
	})

	t.Run("TestDifferentSessions", func(t *testing.T) {
		v := "7f5f965113f9c1b96ef657174bd04433501b220fdecae3128f123311470d4a0ea55dfe6851d17307057c457be59ceb8f76b300d8fc1a0fbe4542c3e16806bc947c286a6ac2a8c2d9face136dbb1f6b8b245c761b2a41e914639ba13711494aac5d1512a756eab9075df3932033ef723aebed294ac9e3699eb10edcd1a62bc0c9a664e3ac21c420c5ea3c143e64ba237e66cc7582fa8d2d19a8a4613e280bf0"
		e, err := decrypt(ctx, v)
		assert.Equal(t, nil, err)
		assert.Equal(t, val, e)
	})
}

func TestBigTest(t *testing.T) {
	key := fmt.Sprint("!;$W=T3rYHXpB'K^")
	ctx := context.WithValue(context.Background(), ContextSecurityKey, key)

	data := []string{
		"AQAe2FljEaseIixwMLbADRq_2xtOQwNNs5K_MmwTYRq32IzwrJ7aqlqRS9DuRsGQoD8ZxUIeKaHL4Er1fDkL7HtY1kzeUckE3FMAZmnF6zBpjf91IjbXYcOorDOnoE-_w5c",
		"AQAe2FljEaseIixwMLbADRq_2xtOQwNNs5K_MmwTYRq32IzwrJ7aqlqRS9DuRsGQoD8ZxUIeKaHL4Er1fDkL7HtY1kzeUckE3FMAZmnF6zBpjf91IjbXYcOorDOnoE-_w4c",
		"AQAe2FljEaseIixwMLbADRq_2xtOQwNNs5K_MmwTYRq32IzwrJ7aqlqRS9DuRsGQoD8ZxUIeKaHL4Er1fDkL7HtY1kzeUckE3FMAZmnF6zBpjf91IjbXYcOorDOnoE-_w3c",
		"AQAe2FljEaseIixwMLbADRq_2xtOQwNNs5K_MmwTYRq32IzwrJ7aqlqRS9DuRsGQoD8ZxUIeKaHL4Er1fDkL7HtY1kzeUckE3FMAZmnF6zBpjf91IjbXYcOorDOnoE-_w2c",
		"AQAe2FljEaseIixwMLbADRq_2xtOQwNNs5K_MmwTYRq32IzwrJ7aqlqRS9DuRsGQoD8ZxUIeKaHL4Er1fDkL7HtY1kzeUckE3FMAZmnF6zBpjf91IjbXYcOorDOnoE-_w1c",
	}

	oldEncs := []string{
		"9da41a8dfee35afbdbc7995d1b8d0855dcf5754d395a5bc137aeb4c060f9d822a1b7c279207134a763cb7e1629659361433091eda2aeea7482f30fc3560109fef3de14f4f146e392b4d32df6e516282fea6890d42cccd13bcc82b278f6dda73e10c464bd321252246adfa0b8c6cfebd3331db2ded6a098925c3b12ac8012819a52a07cfaf647e0e0fa72eef6bafbe7ca175e7aabbfc0a0343b4791e06a0323",
		"ce525a3f3390b8eadb9549a6e9718a1ede35c455acbc12a0faece5c7036871a0157d9742c4ae0d4e69f81f677355f42f35ab040cd66e777283530343fab1cb706fc9b64cef2ebcb74ebdb0fcdc89a4b878a83c8c2f4d59b11d842b2172d95efa285a5b9b700581887e3f80c8f90928d8737f9729d3fe80748bb1d5eea22479cbfca8b885cd727cb37dd4bd1129dc287532efeb1d8dc4b8f21ae807a6a8b2b8",
		"39fbadf339cb7a4ee3400570d1454c1540b0a74176c6296349af55f17847474e22fbc6abc79aacadd5f660d8d6dce691cb90ace4f3e99848d7e314da09b5eda84cb0bd0a8694b285c6ed9d0b7868bb8abf1410cfa27584ad37f91553c11dc127758bb1b6cbba1169db55d34435f7a9bee06e26be31d33f390e8f841fb7b7e655778253fee214cf7b9532c30c8dadf805a8af2f7e347c8e38ea16ed99fb6678",
		"14aa378cbdd69d676b3be85461cb26b8521592d19533ade62da2b6ed86a2b16a575a0d36345df9f7a7f9c4f3d4f35fbfab52b495d8582067527e0e18f20313a47dd370ed9b50f96cf5b82afc08559dd7eef2c4236c59502078998ae0cd861f897ea24e0e67224544d7f4d0051e5469da6c009d4effce99d21b0b35aa64923e456b023de3085a170c9907d08dd8ac1a6db238f31ab91d8470777761d24baf8e",
		"ce14264fbc4cd8333ad256cb4877ab51e94b0ecadea963247fe534826dad461ed29cf0460fad4b6dda76df1524e46459382d343d2abc75f7ea48543d0ab8ad3da204af1e99f6d738409f7b767c9da3a3cc3bb7947687f139dea914560fa1776b1b09ae3f4e3c58ef58c69ee8dfeee099350e4f30174a4f626ab9aae0f7ca4ef8fc3c2d07ba28b3c3d149e6d59b9a96e38cba3a137ed0d523c04c6624655441",
		"faa0451d2d9104e73fa08cfb30103437514bc22d18896ccabc41673a477002ac686f22397eb06c855a9a89bdb72e989b3b4201a17ce1b4393760794805eaf29d3b4689679eb0f12aac6687b80ceccd44dd58f819588ef2e6a43d59d511c47aef96cd0646d62f432277f8dbfe632ac36778a06792c7c14f989aa47869982836c0a4078dc12fd7e23cc8e4b660c779eb0842de2db83bad37ac23b2f0933d59cb",
	}

	for i := 0; i < len(data); i++ {
		enc, err := encrypt(ctx, data[i])
		assert.Equal(t, nil, err)

		enc2, err := encrypt(ctx, data[i])
		assert.Equal(t, nil, err)

		//assert.Equal(t, enc, enc2)

		dec, err := decrypt(ctx, enc)
		assert.Equal(t, nil, err)
		assert.Equal(t, dec, data[i])

		dec2, err := decrypt(ctx, enc2)
		assert.Equal(t, nil, err)
		assert.Equal(t, dec2, data[i])

		//assert.Equal(t, true, false)
	}

	for i := 0; i < len(oldEncs); i++ {
		dec, err := decrypt(ctx, oldEncs[i])
		assert.Equal(t, nil, err)
		found := false
		for j := 0; j < len(data); j++ {
			if dec == data[j] {
				found = true
			}
		}
		assert.Equal(t, true, found)
	}
}
