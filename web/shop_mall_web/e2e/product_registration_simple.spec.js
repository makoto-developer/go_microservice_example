const { test, expect } = require('@playwright/test');

test.describe('商品登録シンプルテスト', () => {
  test('testowner@test.comでログインして商品を登録する', async ({ page }) => {
    // ログイン
    console.log('=== ログイン ===');
    await page.goto('http://localhost:22200/owner/auth');
    await page.fill('input[name="email"]', 'testowner@test.com');
    await page.fill('input[name="password"]', 'secret');
    await page.click('button[type="submit"]');
    await page.waitForTimeout(2000);

    const urlAfterLogin = page.url();
    console.log('ログイン後のURL:', urlAfterLogin);

    if (!urlAfterLogin.includes('/owner/dashboard')) {
      const errorElement = await page.locator('.bg-red-100').count();
      if (errorElement > 0) {
        const errorText = await page.locator('.bg-red-100').textContent();
        throw new Error(`ログイン失敗: ${errorText}`);
      }
      throw new Error(`ダッシュボードにリダイレクトされませんでした: ${urlAfterLogin}`);
    }
    console.log('✅ ログイン成功、ダッシュボードに移動');

    // 商品登録ページに直接アクセス
    console.log('=== 商品登録 ===');
    await page.goto('http://localhost:22200/owner/products/new');
    await page.waitForTimeout(1000);

    // フォームに入力
    const timestamp = Date.now();
    const productName = `テスト商品_${timestamp}`;

    await page.fill('input[name="product[name]"]', productName);
    await page.fill('textarea[name="product[description]"]', 'E2Eテストで登録された商品');
    await page.fill('input[name="product[price]"]', '2980');
    await page.fill('input[name="product[stock_quantity]"]', '50');
    await page.selectOption('select[name="product[category]"]', 'electronics');

    // 公開設定をチェック
    await page.check('input[name="product[published]"]');
    console.log('✅ フォーム入力完了');

    // 登録ボタンをクリック
    await page.click('button[type="submit"]');
    console.log('✅ 登録ボタンクリック');

    // 結果を待機
    await page.waitForTimeout(3000);

    const finalUrl = page.url();
    console.log('最終URL:', finalUrl);

    // 成功またはエラーメッセージを確認
    const successMessage = await page.locator('text=商品を登録しました').count();
    const errorMessage = await page.locator('text=登録に失敗しました').count();

    if (successMessage > 0) {
      console.log('✅✅✅ 商品登録成功！');
      expect(finalUrl).toContain('/owner/products');
      expect(finalUrl).not.toContain('/new');
    } else if (errorMessage > 0) {
      const errorText = await page.locator('.bg-red-100').textContent();
      console.log('❌ エラーメッセージ:', errorText);
      throw new Error(`商品登録失敗: ${errorText}`);
    } else {
      console.log('⚠️ メッセージが表示されませんでした');
      throw new Error('商品登録後、メッセージが表示されませんでした');
    }
  });
});
