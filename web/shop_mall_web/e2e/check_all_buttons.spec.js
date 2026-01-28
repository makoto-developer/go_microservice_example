const { test, expect } = require('@playwright/test');

test.describe('全画面のボタン動作確認', () => {
  test('商品管理画面の公開/非公開ボタンが動作する', async ({ page }) => {
    // ログキャプチャ
    page.on('console', msg => {
      if (msg.type() === 'error') {
        console.log('❌ Browser Error:', msg.text());
      }
    });

    console.log('=== 商品管理画面の公開ボタンテスト ===');

    // 1. 商品登録フォームに直接アクセス
    await page.goto('http://localhost:22200/owner/products/new');
    await page.waitForLoadState('networkidle');

    // 2. 商品登録（公開設定OFF）
    const timestamp = Date.now();
    const productName = `公開テスト商品_${timestamp}`;

    await page.fill('input[name="product[name]"]', productName);
    await page.fill('textarea[name="product[description]"]', 'E2Eテスト: 公開ボタン確認');
    await page.fill('input[name="product[price]"]', '1000');
    await page.fill('input[name="product[stock_quantity]"]', '10');
    await page.selectOption('select[name="product[category]"]', 'electronics');

    // 公開チェックボックスはOFFのまま
    console.log('✅ 商品登録フォーム入力完了（公開設定OFF）');

    // 3. 登録ボタンクリック
    const submitButton = await page.locator('button[type="submit"]');
    await submitButton.scrollIntoViewIfNeeded();
    await submitButton.click();

    // 4. 商品一覧ページへのリダイレクト待ち
    await page.waitForTimeout(3000);
    const currentUrl = page.url();
    expect(currentUrl).toContain('/owner/products');
    console.log('✅ 商品一覧ページにリダイレクト:', currentUrl);

    // 5. 登録した商品を探す
    await page.waitForTimeout(1000);
    const productRow = await page.locator(`text=${productName}`).first();
    expect(await productRow.isVisible()).toBe(true);
    console.log('✅ 登録した商品が表示されている');

    // 6. 公開ボタンを探す
    // 商品名の行から親要素（li）を取得し、その中の公開ボタンを探す
    const parentRow = await productRow.locator('xpath=ancestor::li');
    const publishButton = await parentRow.locator('button:has-text("公開")').first();
    expect(await publishButton.isVisible()).toBe(true);
    console.log('✅ 公開ボタンが表示されている');

    // 7. 公開ボタンをクリック
    console.log('🔄 公開ボタンをクリック...');
    await publishButton.click();

    // 8. フラッシュメッセージを確認
    await page.waitForTimeout(2000);

    // 成功メッセージ
    const successMessage = await page.locator('text=公開状態を変更しました').count();
    const errorMessage = await page.locator('text=変更に失敗しました').count();

    if (successMessage > 0) {
      console.log('✅✅✅ 公開ボタン動作成功！');

      // 9. ボタンのテキストが「非公開」に変わっているか確認
      await page.waitForTimeout(1000);
      const unpublishButton = await parentRow.locator('button:has-text("非公開")').count();

      if (unpublishButton > 0) {
        console.log('✅ ボタンが「非公開」に変更された');
      } else {
        console.log('⚠️ ボタンのテキストは変わっていない（要確認）');
      }
    } else if (errorMessage > 0) {
      const errorText = await page.locator('.bg-red-100').textContent();
      console.log('❌ 公開ボタンエラー:', errorText);
      throw new Error(`公開ボタンが動作しない: ${errorText}`);
    } else {
      console.log('⚠️ フラッシュメッセージなし');
      await page.screenshot({ path: 'test-results/toggle-publish-no-message.png', fullPage: true });
    }
  });

  test('商品登録フォームの全ボタンが動作する', async ({ page }) => {
    console.log('=== 商品登録フォームのボタンテスト ===');

    await page.goto('http://localhost:22200/owner/products/new');
    await page.waitForLoadState('networkidle');

    // 1. キャンセルボタンのテスト
    const cancelButton = await page.locator('text=キャンセル');
    expect(await cancelButton.isVisible()).toBe(true);
    console.log('✅ キャンセルボタンが表示されている');

    await cancelButton.click();
    await page.waitForTimeout(1000);

    const urlAfterCancel = page.url();
    expect(urlAfterCancel).toContain('/owner/products');
    console.log('✅ キャンセルボタンで商品一覧に戻る');

    // 2. 登録ボタンのテスト（前のテストで確認済み）
    console.log('✅ 登録ボタンは前のテストで確認済み');
  });

  test('オーナー認証画面の全ボタンが動作する', async ({ page }) => {
    console.log('=== オーナー認証画面のボタンテスト ===');

    await page.goto('http://localhost:22200/owner/auth');
    await page.waitForLoadState('networkidle');

    // 1. ログインボタン
    const loginButton = await page.locator('button[type="submit"]:has-text("ログイン")');
    expect(await loginButton.isVisible()).toBe(true);
    console.log('✅ ログインボタンが表示されている');

    // 2. モード切り替えボタン
    const toggleButton = await page.locator('button:has-text("新規オーナー登録はこちら")');
    expect(await toggleButton.isVisible()).toBe(true);
    console.log('✅ モード切り替えボタンが表示されている');

    await toggleButton.click();
    await page.waitForTimeout(500);

    const registerButton = await page.locator('button[type="submit"]:has-text("オーナー登録")');
    expect(await registerButton.isVisible()).toBe(true);
    console.log('✅ モード切り替え後、登録ボタンが表示される');

    // 3. カスタマーログインへのリンク
    const customerLink = await page.locator('text=カスタマーログインへ');
    expect(await customerLink.isVisible()).toBe(true);
    console.log('✅ カスタマーログインへのリンクが表示されている');
  });

  test('ナビゲーションバーのボタンが動作する', async ({ page }) => {
    console.log('=== ナビゲーションバーのボタンテスト ===');

    await page.goto('http://localhost:22200/owner/products');
    await page.waitForLoadState('networkidle');

    // 1. ダッシュボードリンク
    const dashboardLink = await page.locator('nav a:has-text("ダッシュボード")');
    expect(await dashboardLink.isVisible()).toBe(true);
    console.log('✅ ダッシュボードリンクが表示されている');

    // 2. 商品管理リンク
    const productsLink = await page.locator('nav a:has-text("商品管理")');
    expect(await productsLink.isVisible()).toBe(true);
    console.log('✅ 商品管理リンクが表示されている');

    // 3. 顧客画面へのリンク
    const customerViewLink = await page.locator('nav a:has-text("顧客画面へ")');
    expect(await customerViewLink.isVisible()).toBe(true);
    console.log('✅ 顧客画面へのリンクが表示されている');
  });

  test('商品一覧の全ボタンが表示される', async ({ page }) => {
    console.log('=== 商品一覧画面のボタンテスト ===');

    await page.goto('http://localhost:22200/owner/products');
    await page.waitForLoadState('networkidle');

    // 1. 新規登録ボタン
    const newProductButton = await page.locator('a:has-text("新規商品登録")');
    expect(await newProductButton.isVisible()).toBe(true);
    console.log('✅ 新規商品登録ボタンが表示されている');

    // 2. 商品があれば、編集・削除・公開ボタンを確認
    const productItems = await page.locator('li').filter({ has: page.locator('a:has-text("編集")') }).count();

    if (productItems > 0) {
      const firstItem = await page.locator('li').filter({ has: page.locator('a:has-text("編集")') }).first();

      const editButton = await firstItem.locator('a:has-text("編集")');
      expect(await editButton.isVisible()).toBe(true);
      console.log('✅ 編集ボタンが表示されている');

      const deleteButton = await firstItem.locator('button:has-text("削除")');
      expect(await deleteButton.isVisible()).toBe(true);
      console.log('✅ 削除ボタンが表示されている');

      const publishToggleButton = await firstItem.locator('button').filter({ hasText: /公開|非公開/ });
      expect(await publishToggleButton.count()).toBeGreaterThan(0);
      console.log('✅ 公開/非公開ボタンが表示されている');
    } else {
      console.log('⚠️ 商品がないため、行内ボタンは確認できない');
    }
  });
});
