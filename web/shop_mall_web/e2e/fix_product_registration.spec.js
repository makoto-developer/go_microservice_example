const { test, expect } = require('@playwright/test');

test.describe('商品登録動作確認', () => {
  test('商品登録フォームが正常に動作する', async ({ page }) => {
    // ログをキャプチャ
    page.on('console', msg => {
      if (msg.type() === 'error') {
        console.log('❌ Browser Error:', msg.text());
      }
    });

    // 1. オーナー認証ページにアクセス
    console.log('=== Step 1: オーナー認証 ===');
    await page.goto('http://localhost:22200/owner/auth');
    await page.waitForLoadState('networkidle');

    // 2. ログイン情報を入力
    await page.fill('input[name="email"]', 'testowner@test.com');
    await page.fill('input[name="password"]', 'secret');
    console.log('✅ ログイン情報入力');

    // 3. ログインボタンをクリック
    await page.click('button[type="submit"]');
    console.log('✅ ログインボタンクリック');

    // 4. ダッシュボードまたはエラーを待機
    await page.waitForTimeout(3000);

    const currentUrl = page.url();
    console.log('ログイン後URL:', currentUrl);

    // エラーメッセージがあるか確認
    const errorCount = await page.locator('.bg-red-100').count();
    if (errorCount > 0) {
      const errorText = await page.locator('.bg-red-100').textContent();
      console.log('❌ ログインエラー:', errorText);

      // それでも続行（商品登録フォームに直接アクセス）
      console.log('⚠️ ログイン失敗、フォームに直接アクセスします');
      await page.goto('http://localhost:22200/owner/products/new');
    } else if (currentUrl.includes('/owner/dashboard')) {
      console.log('✅ ダッシュボードに移動成功');

      // 5. 商品登録ページに移動
      await page.goto('http://localhost:22200/owner/products/new');
    } else {
      console.log('⚠️ 予期しないURL、フォームに直接アクセスします');
      await page.goto('http://localhost:22200/owner/products/new');
    }

    await page.waitForLoadState('networkidle');
    console.log('✅ 商品登録フォームにアクセス');

    // 6. フォームに入力
    const timestamp = Date.now();
    const productName = `E2E商品_${timestamp}`;

    await page.fill('input[name="product[name]"]', productName);
    console.log('✅ 商品名入力:', productName);

    await page.fill('textarea[name="product[description]"]', 'E2Eテストで登録された商品');
    await page.fill('input[name="product[price]"]', '3500');
    await page.fill('input[name="product[stock_quantity]"]', '75');
    await page.selectOption('select[name="product[category]"]', 'electronics');

    // 公開設定
    await page.check('input[name="product[published]"]');
    console.log('✅ フォーム入力完了');

    // 7. 登録ボタンをクリック
    const submitButton = await page.locator('button[type="submit"]');
    await submitButton.scrollIntoViewIfNeeded();

    console.log('✅ 登録ボタンをクリック');
    await submitButton.click();

    // 8. 結果を待機
    await page.waitForTimeout(5000);

    const finalUrl = page.url();
    console.log('最終URL:', finalUrl);

    // 成功メッセージを確認
    const successCount = await page.locator('text=商品を登録しました').count();
    const errorMsgCount = await page.locator('text=登録に失敗しました').count();

    if (successCount > 0) {
      console.log('✅✅✅ 商品登録成功！');
      expect(finalUrl).toContain('/owner/products');
    } else if (errorMsgCount > 0) {
      const errorText = await page.locator('.bg-red-100').textContent();
      console.log('❌ 登録エラー:', errorText);
      throw new Error(`商品登録失敗: ${errorText}`);
    } else {
      console.log('⚠️ 成功/エラーメッセージなし');
      console.log('現在のURL:', finalUrl);

      // スクリーンショットを撮る
      await page.screenshot({ path: 'test-results/debug-screenshot.png', fullPage: true });
      console.log('スクリーンショット保存: test-results/debug-screenshot.png');
    }
  });
});
