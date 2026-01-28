const { test, expect } = require('@playwright/test');

test.describe('商品登録フルフロー', () => {
  test('オーナー登録から商品登録までの完全なフロー', async ({ page }) => {
    const timestamp = Date.now();
    const testEmail = `owner_${timestamp}@test.com`;
    const testPassword = 'password123';

    // === オーナー登録 ===
    console.log('=== Step 1: オーナー登録 ===');
    await page.goto('http://localhost:22200/owner/auth');

    // 新規登録モードに切り替え
    await page.click('button:has-text("新規オーナー登録はこちら")');
    await page.waitForTimeout(500);
    console.log('✅ 新規登録モードに切り替え');

    // オーナー登録
    await page.fill('input[name="email"]', testEmail);
    await page.fill('input[name="password"]', testPassword);
    await page.click('button[type="submit"]');
    console.log(`✅ オーナー登録実行: ${testEmail}`);

    // ショップ登録ページにリダイレクトされるまで待機
    await page.waitForTimeout(2000);
    const urlAfterRegister = page.url();
    console.log('登録後のURL:', urlAfterRegister);

    if (!urlAfterRegister.includes('/owner/shop/register')) {
      const errorElement = await page.locator('.bg-red-100').count();
      if (errorElement > 0) {
        const errorText = await page.locator('.bg-red-100').textContent();
        throw new Error(`オーナー登録失敗: ${errorText}`);
      }
      throw new Error(`ショップ登録ページにリダイレクトされませんでした: ${urlAfterRegister}`);
    }
    console.log('✅ ショップ登録ページにリダイレクト成功');

    // === ショップ登録 ===
    console.log('=== Step 2: ショップ登録 ===');
    await page.fill('input[name="shop[name]"]', `テストショップ_${timestamp}`);
    await page.fill('textarea[name="shop[description]"]', 'E2Eテスト用のショップです');
    await page.fill('input[name="shop[owner_name]"]', 'テストオーナー');
    await page.fill('input[name="shop[phone_number]"]', '090-1234-5678');
    await page.fill('input[name="shop[business_hours]"]', '9:00-18:00');
    await page.fill('textarea[name="shop[return_policy]"]', '商品到着後7日以内');
    console.log('✅ ショップ情報入力完了');

    await page.click('button[type="submit"]');
    console.log('✅ ショップ登録実行');

    // ダッシュボードにリダイレクトされるまで待機
    await page.waitForTimeout(3000);
    const urlAfterShop = page.url();
    console.log('ショップ登録後のURL:', urlAfterShop);

    if (!urlAfterShop.includes('/owner/dashboard')) {
      const errorElement = await page.locator('.bg-red-100').count();
      if (errorElement > 0) {
        const errorText = await page.locator('.bg-red-100').textContent();
        console.log('❌ エラー:', errorText);
      }
      throw new Error(`ダッシュボードにリダイレクトされませんでした: ${urlAfterShop}`);
    }
    console.log('✅ ダッシュボードに移動');

    // === 商品登録 ===
    console.log('=== Step 3: 商品登録 ===');

    // 商品管理ページへ移動
    await page.click('a[href="/owner/products"]');
    await page.waitForTimeout(1000);
    console.log('✅ 商品管理ページへ移動');

    // 新規商品登録
    await page.click('a[href="/owner/products/new"]');
    await page.waitForTimeout(1000);
    console.log('✅ 商品登録フォームへ移動');

    // フォームに入力
    const productName = `テスト商品_${timestamp}`;
    await page.fill('input[name="product[name]"]', productName);
    await page.fill('textarea[name="product[description]"]', 'E2Eテストで登録された商品です');
    await page.fill('input[name="product[price]"]', '1500');
    await page.fill('input[name="product[stock_quantity]"]', '100');
    await page.selectOption('select[name="product[category]"]', 'electronics');

    // 公開設定をチェック
    await page.check('input[name="product[published]"]');
    console.log('✅ フォーム入力完了（公開設定: ON）');

    // 登録ボタンをクリック
    await page.click('button[type="submit"]');
    console.log('✅ 登録ボタンクリック');

    // 結果を待機
    await page.waitForTimeout(3000);

    const finalUrl = page.url();
    console.log('最終URL:', finalUrl);

    // 成功メッセージまたはエラーメッセージを確認
    const successMessage = await page.locator('text=商品を登録しました').count();
    const errorMessage = await page.locator('text=登録に失敗しました').count();

    if (successMessage > 0) {
      console.log('✅ 成功メッセージ表示');
      expect(finalUrl).toContain('/owner/products');
      expect(finalUrl).not.toContain('/new');
      console.log('✅ 商品一覧ページにリダイレクト成功');
    } else if (errorMessage > 0) {
      const errorText = await page.locator('.bg-red-100').textContent();
      console.log('❌ エラーメッセージ:', errorText);
      throw new Error(`商品登録失敗: ${errorText}`);
    } else {
      console.log('⚠️ 成功/エラーメッセージが表示されませんでした');
      throw new Error('商品登録後、メッセージが表示されませんでした');
    }

    console.log('=== 全ステップ完了 ===');
  });
});
