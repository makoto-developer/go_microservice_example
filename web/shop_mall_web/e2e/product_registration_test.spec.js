const { test, expect } = require('@playwright/test');

test.describe('商品登録実際のテスト', () => {
  test('実際に商品を登録できる', async ({ page }) => {
    // オーナー認証画面にアクセス
    await page.goto('http://localhost:22200/owner/auth');
    console.log('✅ オーナー認証画面にアクセス');

    // ページの状態を確認
    await page.waitForSelector('input[name="email"]');
    console.log('✅ メールアドレス入力欄が表示されました');

    // ログイン情報を入力
    await page.fill('input[name="email"]', 'owner1@example.com');
    console.log('✅ メールアドレスを入力: owner1@example.com');

    await page.fill('input[name="password"]', 'password123');
    console.log('✅ パスワードを入力');

    // 送信ボタンを探す
    const submitButton = await page.locator('button[type="submit"]');
    const buttonText = await submitButton.textContent();
    console.log('送信ボタンのテキスト:', buttonText);

    // ログインボタンをクリック
    await submitButton.click();
    console.log('✅ ログインボタンをクリック');

    // 少し待機
    await page.waitForTimeout(2000);

    // 現在のURLを確認
    const urlAfterLogin = page.url();
    console.log('ログイン後のURL:', urlAfterLogin);

    // エラーメッセージがあるか確認
    const errorElement = await page.locator('.bg-red-100').count();
    if (errorElement > 0) {
      const errorText = await page.locator('.bg-red-100').textContent();
      console.log('❌ エラーメッセージ:', errorText);
      throw new Error(`ログイン失敗: ${errorText}`);
    }

    // ダッシュボードへのリダイレクトを確認（タイムアウトなし）
    if (urlAfterLogin.includes('/owner/dashboard')) {
      console.log('✅ ダッシュボードにリダイレクトされました');
    } else {
      console.log('⚠️ ダッシュボードにリダイレクトされていません:', urlAfterLogin);
      throw new Error(`ログイン後、期待されるページにリダイレクトされませんでした: ${urlAfterLogin}`);
    }

    // 商品管理ページへ移動
    await page.click('a[href="/owner/products"]');
    await page.waitForTimeout(1000);
    console.log('✅ 商品管理リンクをクリック');
    console.log('現在のURL:', page.url());

    // 新規商品登録ページへ移動
    const newProductLink = await page.locator('a[href="/owner/products/new"]');
    const linkCount = await newProductLink.count();
    console.log('新規商品登録リンクの数:', linkCount);

    if (linkCount > 0) {
      await newProductLink.click();
      await page.waitForTimeout(1000);
      console.log('✅ 商品登録フォームへ移動');
      console.log('現在のURL:', page.url());
    } else {
      console.log('❌ 新規商品登録リンクが見つかりません');
      throw new Error('新規商品登録リンクが見つかりません');
    }

    // フォームに入力
    const timestamp = Date.now();
    const productName = `テスト商品_${timestamp}`;

    await page.fill('input[name="product[name]"]', productName);
    await page.fill('textarea[name="product[description]"]', 'これはE2Eテストで登録された商品です');
    await page.fill('input[name="product[price]"]', '1500');
    await page.fill('input[name="product[stock_quantity]"]', '100');
    await page.selectOption('select[name="product[category]"]', 'electronics');
    console.log('✅ フォーム入力完了');

    // 登録ボタンをクリック
    const registerButton = await page.locator('button[type="submit"]');
    const registerButtonText = await registerButton.textContent();
    console.log('登録ボタンのテキスト:', registerButtonText);

    await registerButton.click();
    console.log('✅ 登録ボタンをクリック');

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
    } else if (errorMessage > 0) {
      const errorText = await page.locator('.bg-red-100').textContent();
      console.log('❌ エラーメッセージ:', errorText);
      throw new Error(`商品登録失敗: ${errorText}`);
    } else {
      console.log('⚠️ 成功/エラーメッセージが表示されませんでした');
      throw new Error('商品登録後、メッセージが表示されませんでした');
    }
  });
});
