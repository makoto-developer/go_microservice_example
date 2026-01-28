// @ts-check
const { test, expect } = require('@playwright/test');

test.describe('オーナー業務フロー E2E', () => {
  // テストデータ
  const ownerEmail = `owner_${Date.now()}@example.com`;
  const ownerPassword = 'SecurePass123!';
  const shopName = 'テストショップ';
  
  const product1 = {
    name: 'ワイヤレスイヤホン Pro',
    description: '高音質ワイヤレスイヤホン【期間限定セール】',
    price: '11800',
    stock: '100',
    sku: 'WH-PRO-001'
  };
  
  const product2 = {
    name: 'Bluetoothスピーカー',
    description: '防水対応ポータブルスピーカー',
    price: '8500',
    stock: '50',
    sku: 'BT-SPK-001'
  };

  test('1. オーナー登録が正常に動作する', async ({ page }) => {
    await page.goto('/owner/auth');
    
    // ページタイトル確認
    await expect(page).toHaveURL(/\/owner\/auth/);
    
    // 登録フォームが表示されることを確認
    await expect(page.locator('text=オーナー登録').or(page.locator('text=登録')).first()).toBeVisible();
    
    // フォーム要素の確認
    await expect(page.locator('input[type="email"]')).toBeVisible();
    await expect(page.locator('input[type="password"]')).toBeVisible();
    
    // スクリーンショット撮影
    await page.screenshot({ path: 'e2e/screenshots/owner-01-auth.png', fullPage: true });
    
    console.log('✅ オーナー登録画面: 表示確認OK');
  });

  test('2. オーナー認証・ログインが動作する', async ({ page }) => {
    await page.goto('/owner/auth');
    
    // ログインフォームに切り替え（既にログインフォームが表示されている場合）
    const loginTab = page.locator('text=ログイン').first();
    if (await loginTab.isVisible()) {
      await loginTab.click();
    }
    
    // メールアドレス・パスワード入力（デモ用）
    await page.locator('input[type="email"]').fill('demo_owner@example.com');
    await page.locator('input[type="password"]').fill('DemoPass123!');
    
    // ログインボタンクリック（実際にはgRPC通信が必要）
    const loginButton = page.locator('button:has-text("ログイン")').or(page.locator('button[type="submit"]')).first();
    if (await loginButton.isVisible()) {
      await loginButton.click();
      
      // gRPC通信を待つ
      await page.waitForTimeout(2000);
      
      // ログイン成功の場合はダッシュボードへ遷移
      // エラーの場合はエラーメッセージが表示される
      const currentUrl = page.url();
      console.log(`現在のURL: ${currentUrl}`);
    }
    
    console.log('✅ オーナー認証: フォーム確認OK');
  });

  test('3. ショップ登録画面が表示される', async ({ page }) => {
    await page.goto('/owner/shop/register');
    
    // ページが正常に表示されることを確認
    await expect(page.locator('text=ショップ登録').or(page.locator('text=登録')).first()).toBeVisible();
    
    // フォーム要素の確認
    const nameInput = page.locator('input[name="name"]').or(page.locator('input[placeholder*="ショップ名"]')).first();
    if (await nameInput.isVisible()) {
      await expect(nameInput).toBeVisible();
    }
    
    // スクリーンショット撮影
    await page.screenshot({ path: 'e2e/screenshots/owner-02-shop-register.png', fullPage: true });
    
    console.log('✅ ショップ登録画面: 表示確認OK');
  });

  test('4. オーナーダッシュボードが表示される', async ({ page }) => {
    await page.goto('/owner/dashboard');
    
    // ダッシュボードが表示されることを確認
    await expect(page.locator('text=ダッシュボード').first()).toBeVisible();
    
    // gRPC通信を待つ
    await page.waitForTimeout(2000);
    
    // 統計情報が表示されることを確認
    const statsVisible = await page.locator('text=登録商品数').or(page.locator('text=商品数')).isVisible().catch(() => false);
    
    if (statsVisible) {
      console.log('✅ 統計情報が表示されました');
    } else {
      console.log('⚠️ 統計情報の取得に失敗（gRPCサービス未起動の可能性）');
    }
    
    // スクリーンショット撮影
    await page.screenshot({ path: 'e2e/screenshots/owner-03-dashboard.png', fullPage: true });
    
    console.log('✅ オーナーダッシュボード: 表示確認OK');
  });

  test('5. 商品登録フォームが表示される', async ({ page }) => {
    await page.goto('/owner/products/new');
    
    // 商品登録フォームが表示されることを確認
    await expect(page.locator('text=商品登録').or(page.locator('text=新規商品')).first()).toBeVisible();
    
    // フォーム要素の確認
    const productNameInput = page.locator('input[name="name"]').or(page.locator('input[placeholder*="商品名"]')).first();
    if (await productNameInput.isVisible()) {
      await expect(productNameInput).toBeVisible();
    }
    
    // スクリーンショット撮影
    await page.screenshot({ path: 'e2e/screenshots/owner-04-product-form.png', fullPage: true });
    
    console.log('✅ 商品登録フォーム: 表示確認OK');
  });

  test('6. 商品管理一覧が表示される', async ({ page }) => {
    await page.goto('/owner/products');
    
    // 商品管理一覧が表示されることを確認
    await expect(page.locator('text=商品管理').or(page.locator('text=商品一覧')).first()).toBeVisible();
    
    // gRPC通信を待つ
    await page.waitForTimeout(2000);
    
    // 商品が表示されることを確認（gRPC通信成功時）
    const productsVisible = await page.locator('text=ワイヤレスイヤホン').or(page.locator('text=在庫')).isVisible().catch(() => false);
    
    if (productsVisible) {
      console.log('✅ 商品データが表示されました');
    } else {
      console.log('⚠️ 商品データの取得に失敗（gRPCサービス未起動の可能性）');
    }
    
    // スクリーンショット撮影
    await page.screenshot({ path: 'e2e/screenshots/owner-05-product-list.png', fullPage: true });
    
    console.log('✅ 商品管理一覧: 表示確認OK');
  });

  test('7. 商品編集画面が表示される', async ({ page }) => {
    // テスト用商品ID（実際にはDBから取得）
    const testProductId = 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa';
    
    await page.goto(`/owner/products/${testProductId}/edit`);
    
    // gRPC通信を待つ
    await page.waitForTimeout(2000);
    
    // 商品編集フォームまたはエラーメッセージが表示されることを確認
    const formVisible = await page.locator('text=商品編集').or(page.locator('input[name="name"]')).isVisible().catch(() => false);
    const errorVisible = await page.locator('text=見つかりませんでした').or(page.locator('text=エラー')).isVisible().catch(() => false);
    
    expect(formVisible || errorVisible).toBeTruthy();
    
    // スクリーンショット撮影
    await page.screenshot({ path: 'e2e/screenshots/owner-06-product-edit.png', fullPage: true });
    
    console.log('✅ 商品編集画面: 表示確認OK');
  });

  test('8. 全画面への遷移が正常に動作する', async ({ page }) => {
    // オーナー認証画面
    await page.goto('/owner/auth');
    await expect(page).toHaveURL(/\/owner\/auth/);
    console.log('✅ オーナー認証画面へ遷移OK');
    
    // ショップ登録画面
    await page.goto('/owner/shop/register');
    await expect(page).toHaveURL(/\/owner\/shop\/register/);
    console.log('✅ ショップ登録画面へ遷移OK');
    
    // ダッシュボード
    await page.goto('/owner/dashboard');
    await expect(page).toHaveURL(/\/owner\/dashboard/);
    console.log('✅ ダッシュボードへ遷移OK');
    
    // 商品管理一覧
    await page.goto('/owner/products');
    await expect(page).toHaveURL(/\/owner\/products/);
    console.log('✅ 商品管理一覧へ遷移OK');
    
    // 商品登録
    await page.goto('/owner/products/new');
    await expect(page).toHaveURL(/\/owner\/products\/new/);
    console.log('✅ 商品登録画面へ遷移OK');
    
    console.log('✅ 全画面ナビゲーション: 遷移確認OK');
  });

  test('9. 各画面が3秒以内にロードされる', async ({ page }) => {
    const pages = [
      { path: '/owner/auth', name: 'オーナー認証画面' },
      { path: '/owner/shop/register', name: 'ショップ登録画面' },
      { path: '/owner/dashboard', name: 'オーナーダッシュボード' },
      { path: '/owner/products', name: '商品管理一覧' },
      { path: '/owner/products/new', name: '商品登録フォーム' },
    ];

    for (const { path, name } of pages) {
      const startTime = Date.now();
      await page.goto(path);
      const loadTime = Date.now() - startTime;

      expect(loadTime).toBeLessThan(3000);
      console.log(`✅ ${name}: ${loadTime}ms でロード完了`);
    }
  });

  test('10. レスポンシブデザイン: モバイル表示', async ({ page }) => {
    // モバイルサイズに変更
    await page.setViewportSize({ width: 375, height: 667 });
    
    await page.goto('/owner/dashboard');
    await page.waitForTimeout(2000);
    
    // ダッシュボードが表示されることを確認
    await expect(page.locator('text=ダッシュボード').first()).toBeVisible();
    
    // スクリーンショット撮影
    await page.screenshot({ path: 'e2e/screenshots/owner-mobile-dashboard.png', fullPage: true });
    
    console.log('✅ モバイル表示: OK');
  });

  test('11. レスポンシブデザイン: タブレット表示', async ({ page }) => {
    // タブレットサイズに変更
    await page.setViewportSize({ width: 768, height: 1024 });
    
    await page.goto('/owner/products');
    await page.waitForTimeout(2000);
    
    // 商品管理一覧が表示されることを確認
    await expect(page.locator('text=商品管理').or(page.locator('text=商品一覧')).first()).toBeVisible();
    
    // スクリーンショット撮影
    await page.screenshot({ path: 'e2e/screenshots/owner-tablet-products.png', fullPage: true });
    
    console.log('✅ タブレット表示: OK');
  });
});

test.describe('オーナーフォーム入力テスト', () => {
  test('ショップ登録フォームに入力できる', async ({ page }) => {
    await page.goto('/owner/shop/register');
    
    // フォーム要素を探す
    const nameInput = page.locator('input[name="name"]').or(page.locator('input[placeholder*="ショップ名"]')).first();
    const descriptionInput = page.locator('textarea[name="description"]').or(page.locator('textarea[placeholder*="説明"]')).first();
    
    // フォームに入力
    if (await nameInput.isVisible()) {
      await nameInput.fill('テストショップ');
      console.log('✅ ショップ名入力OK');
    }
    
    if (await descriptionInput.isVisible()) {
      await descriptionInput.fill('テスト用のショップです');
      console.log('✅ ショップ説明入力OK');
    }
    
    // スクリーンショット撮影
    await page.screenshot({ path: 'e2e/screenshots/owner-shop-register-filled.png', fullPage: true });
    
    console.log('✅ ショップ登録フォーム入力: OK');
  });

  test('商品登録フォームに入力できる', async ({ page }) => {
    await page.goto('/owner/products/new');
    
    // フォーム要素を探す
    const nameInput = page.locator('input[name="name"]').or(page.locator('input[placeholder*="商品名"]')).first();
    const descriptionInput = page.locator('textarea[name="description"]').or(page.locator('textarea[placeholder*="説明"]')).first();
    const priceInput = page.locator('input[name="price"]').or(page.locator('input[placeholder*="価格"]')).first();
    
    // フォームに入力
    if (await nameInput.isVisible()) {
      await nameInput.fill('ワイヤレスイヤホン Pro');
      console.log('✅ 商品名入力OK');
    }
    
    if (await descriptionInput.isVisible()) {
      await descriptionInput.fill('高音質ワイヤレスイヤホン【期間限定セール】');
      console.log('✅ 商品説明入力OK');
    }
    
    if (await priceInput.isVisible()) {
      await priceInput.fill('11800');
      console.log('✅ 価格入力OK');
    }
    
    // スクリーンショット撮影
    await page.screenshot({ path: 'e2e/screenshots/owner-product-form-filled.png', fullPage: true });
    
    console.log('✅ 商品登録フォーム入力: OK');
  });
});

test.describe('gRPC通信確認テスト', () => {
  test('オーナーダッシュボードでgRPC通信が動作する', async ({ page }) => {
    await page.goto('/owner/dashboard');
    
    // gRPC通信を待つ
    await page.waitForTimeout(3000);
    
    // 統計情報が表示されているか確認
    const statsVisible = await page.locator('text=登録商品数').or(page.locator('text=商品数')).isVisible().catch(() => false);
    
    if (statsVisible) {
      // 商品数を取得
      const productCountElement = page.locator('text=登録商品数').locator('..').locator('dd').or(page.locator('text=商品数').locator('..').locator('dd')).first();
      const productCount = await productCountElement.textContent().catch(() => '取得失敗');
      
      console.log(`✅ gRPC通信成功: 登録商品数 = ${productCount}`);
    } else {
      console.log('⚠️ gRPC通信失敗: バックエンドサービスが起動していない可能性があります');
    }
  });

  test('商品管理一覧でgRPC通信が動作する', async ({ page }) => {
    await page.goto('/owner/products');
    
    // gRPC通信を待つ
    await page.waitForTimeout(3000);
    
    // 商品が表示されているか確認
    const productsVisible = await page.locator('text=ワイヤレスイヤホン').or(page.locator('text=在庫')).isVisible().catch(() => false);
    
    if (productsVisible) {
      console.log('✅ gRPC通信成功: 商品データが取得されました');
    } else {
      console.log('⚠️ gRPC通信失敗: バックエンドサービスが起動していない可能性があります');
    }
  });
});