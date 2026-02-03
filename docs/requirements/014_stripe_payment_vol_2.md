# Stripe応用機能実装ガイド
## 複数の支払い方法・返金処理・請求書自動生成

---

## 目次

1. [複数の支払い方法のサポート](#複数の支払い方法のサポート)
   - Payment Intentの使用
   - 複数の支払い方法の追加
   - デフォルト支払い方法の設定
2. [返金処理の実装](#返金処理の実装)
   - 全額返金
   - 部分返金
   - 返金履歴の管理
3. [請求書の自動生成](#請求書の自動生成)
   - 請求書の作成
   - PDF生成
   - メール送信
4. [実践的な実装例](#実践的な実装例)

---

## 複数の支払い方法のサポート

### 1.1 Payment Intentを使用した実装

Checkout Sessionではなく、Payment Intentを使用することで、複数の支払い方法をサポートできます。

#### バックエンド実装（Go）

```go
package main

import (
    "net/http"
    "github.com/labstack/echo/v4"
    "github.com/stripe/stripe-go/v76"
    "github.com/stripe/stripe-go/v76/paymentintent"
    "github.com/stripe/stripe-go/v76/paymentmethod"
    "github.com/stripe/stripe-go/v76/customer"
)

// Payment Intent作成
type CreatePaymentIntentRequest struct {
    Amount            int64    `json:"amount"`
    Currency          string   `json:"currency"`
    PaymentMethodTypes []string `json:"paymentMethodTypes"`
    CustomerID        string   `json:"customerId,omitempty"`
}

func createPaymentIntent(c echo.Context) error {
    var req CreatePaymentIntentRequest
    if err := c.Bind(&req); err != nil {
        return c.JSON(http.StatusBadRequest, map[string]string{
            "error": err.Error(),
        })
    }

    // デフォルトの支払い方法タイプ
    if len(req.PaymentMethodTypes) == 0 {
        req.PaymentMethodTypes = []string{"card", "konbini", "paypay"}
    }

    params := &stripe.PaymentIntentParams{
        Amount:   stripe.Int64(req.Amount),
        Currency: stripe.String(req.Currency),
        PaymentMethodTypes: stripe.StringSlice(req.PaymentMethodTypes),
        AutomaticPaymentMethods: &stripe.PaymentIntentAutomaticPaymentMethodsParams{
            Enabled: stripe.Bool(true),
        },
    }

    // 顧客IDがある場合は設定
    if req.CustomerID != "" {
        params.Customer = stripe.String(req.CustomerID)
    }

    pi, err := paymentintent.New(params)
    if err != nil {
        return c.JSON(http.StatusInternalServerError, map[string]string{
            "error": err.Error(),
        })
    }

    return c.JSON(http.StatusOK, map[string]string{
        "clientSecret": pi.ClientSecret,
        "paymentIntentId": pi.ID,
    })
}

// 顧客の作成
type CreateCustomerRequest struct {
    Email string `json:"email"`
    Name  string `json:"name"`
}

func createCustomer(c echo.Context) error {
    var req CreateCustomerRequest
    if err := c.Bind(&req); err != nil {
        return c.JSON(http.StatusBadRequest, map[string]string{
            "error": err.Error(),
        })
    }

    params := &stripe.CustomerParams{
        Email: stripe.String(req.Email),
        Name:  stripe.String(req.Name),
    }

    cust, err := customer.New(params)
    if err != nil {
        return c.JSON(http.StatusInternalServerError, map[string]string{
            "error": err.Error(),
        })
    }

    return c.JSON(http.StatusOK, map[string]string{
        "customerId": cust.ID,
    })
}
```

#### フロントエンド実装（Next.js）

```typescript
'use client';

import { useState, useEffect } from 'react';
import { loadStripe, Stripe, StripeElements } from '@stripe/stripe-js';
import {
  Elements,
  PaymentElement,
  useStripe,
  useElements,
} from '@stripe/react-stripe-js';

const stripePromise = loadStripe(
  process.env.NEXT_PUBLIC_STRIPE_PUBLISHABLE_KEY!
);

// チェックアウトフォームコンポーネント
function CheckoutForm({ clientSecret }: { clientSecret: string }) {
  const stripe = useStripe();
  const elements = useElements();
  const [message, setMessage] = useState<string>('');
  const [isLoading, setIsLoading] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    if (!stripe || !elements) {
      return;
    }

    setIsLoading(true);

    const { error } = await stripe.confirmPayment({
      elements,
      confirmParams: {
        return_url: `${window.location.origin}/payment-complete`,
      },
    });

    if (error) {
      setMessage(error.message || '支払いに失敗しました');
    }

    setIsLoading(false);
  };

  return (
    <form onSubmit={handleSubmit} className="max-w-md mx-auto p-6">
      <PaymentElement />
      <button
        disabled={isLoading || !stripe || !elements}
        className="w-full mt-4 bg-blue-600 text-white py-3 rounded-lg disabled:opacity-50"
      >
        {isLoading ? '処理中...' : '支払う'}
      </button>
      {message && <div className="mt-4 text-red-600">{message}</div>}
    </form>
  );
}

// メインコンポーネント
export default function PaymentPage() {
  const [clientSecret, setClientSecret] = useState('');

  useEffect(() => {
    // Payment Intent作成
    fetch('http://localhost:8080/create-payment-intent', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        amount: 2000,
        currency: 'jpy',
        paymentMethodTypes: ['card', 'konbini', 'paypay'],
      }),
    })
      .then((res) => res.json())
      .then((data) => setClientSecret(data.clientSecret));
  }, []);

  const options = {
    clientSecret,
    appearance: {
      theme: 'stripe' as const,
    },
  };

  return (
    <div className="min-h-screen bg-gray-50 py-12">
      <h1 className="text-3xl font-bold text-center mb-8">
        お支払い方法を選択
      </h1>
      {clientSecret && (
        <Elements options={options} stripe={stripePromise}>
          <CheckoutForm clientSecret={clientSecret} />
        </Elements>
      )}
    </div>
  );
}
```

### 1.2 保存済み支払い方法の管理

#### 支払い方法の追加

```go
// 支払い方法の追加
type AttachPaymentMethodRequest struct {
    PaymentMethodID string `json:"paymentMethodId"`
    CustomerID      string `json:"customerId"`
}

func attachPaymentMethod(c echo.Context) error {
    var req AttachPaymentMethodRequest
    if err := c.Bind(&req); err != nil {
        return c.JSON(http.StatusBadRequest, map[string]string{
            "error": err.Error(),
        })
    }

    params := &stripe.PaymentMethodAttachParams{
        Customer: stripe.String(req.CustomerID),
    }

    pm, err := paymentmethod.Attach(req.PaymentMethodID, params)
    if err != nil {
        return c.JSON(http.StatusInternalServerError, map[string]string{
            "error": err.Error(),
        })
    }

    return c.JSON(http.StatusOK, map[string]interface{}{
        "paymentMethod": pm,
    })
}

// 顧客の支払い方法一覧取得
func listPaymentMethods(c echo.Context) error {
    customerID := c.QueryParam("customerId")
    
    params := &stripe.PaymentMethodListParams{
        Customer: stripe.String(customerID),
        Type:     stripe.String("card"),
    }

    i := paymentmethod.List(params)
    var methods []*stripe.PaymentMethod
    
    for i.Next() {
        methods = append(methods, i.PaymentMethod())
    }

    if err := i.Err(); err != nil {
        return c.JSON(http.StatusInternalServerError, map[string]string{
            "error": err.Error(),
        })
    }

    return c.JSON(http.StatusOK, map[string]interface{}{
        "paymentMethods": methods,
    })
}

// デフォルト支払い方法の設定
type SetDefaultPaymentMethodRequest struct {
    CustomerID      string `json:"customerId"`
    PaymentMethodID string `json:"paymentMethodId"`
}

func setDefaultPaymentMethod(c echo.Context) error {
    var req SetDefaultPaymentMethodRequest
    if err := c.Bind(&req); err != nil {
        return c.JSON(http.StatusBadRequest, map[string]string{
            "error": err.Error(),
        })
    }

    params := &stripe.CustomerParams{
        InvoiceSettings: &stripe.CustomerInvoiceSettingsParams{
            DefaultPaymentMethod: stripe.String(req.PaymentMethodID),
        },
    }

    cust, err := customer.Update(req.CustomerID, params)
    if err != nil {
        return c.JSON(http.StatusInternalServerError, map[string]string{
            "error": err.Error(),
        })
    }

    return c.JSON(http.StatusOK, map[string]string{
        "status": "success",
        "customerId": cust.ID,
    })
}
```

#### フロントエンド（保存済み支払い方法の表示）

```typescript
'use client';

import { useState, useEffect } from 'react';

interface PaymentMethod {
  id: string;
  card: {
    brand: string;
    last4: string;
    exp_month: number;
    exp_year: number;
  };
}

export default function SavedPaymentMethods({ customerId }: { customerId: string }) {
  const [methods, setMethods] = useState<PaymentMethod[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    fetch(`http://localhost:8080/payment-methods?customerId=${customerId}`)
      .then((res) => res.json())
      .then((data) => {
        setMethods(data.paymentMethods);
        setLoading(false);
      });
  }, [customerId]);

  const handleSetDefault = async (paymentMethodId: string) => {
    await fetch('http://localhost:8080/set-default-payment-method', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        customerId,
        paymentMethodId,
      }),
    });
  };

  if (loading) return <div>読み込み中...</div>;

  return (
    <div className="max-w-2xl mx-auto p-6">
      <h2 className="text-2xl font-bold mb-4">保存済みの支払い方法</h2>
      <div className="space-y-4">
        {methods.map((method) => (
          <div
            key={method.id}
            className="border rounded-lg p-4 flex justify-between items-center"
          >
            <div>
              <p className="font-semibold">
                {method.card.brand.toUpperCase()} •••• {method.card.last4}
              </p>
              <p className="text-sm text-gray-600">
                有効期限: {method.card.exp_month}/{method.card.exp_year}
              </p>
            </div>
            <button
              onClick={() => handleSetDefault(method.id)}
              className="bg-blue-600 text-white px-4 py-2 rounded"
            >
              デフォルトに設定
            </button>
          </div>
        ))}
      </div>
    </div>
  );
}
```

---

## 返金処理の実装

### 2.1 全額返金

```go
package main

import (
    "net/http"
    "github.com/labstack/echo/v4"
    "github.com/stripe/stripe-go/v76"
    "github.com/stripe/stripe-go/v76/refund"
)

// 全額返金
type CreateRefundRequest struct {
    PaymentIntentID string `json:"paymentIntentId"`
    Reason          string `json:"reason,omitempty"`
}

func createRefund(c echo.Context) error {
    var req CreateRefundRequest
    if err := c.Bind(&req); err != nil {
        return c.JSON(http.StatusBadRequest, map[string]string{
            "error": err.Error(),
        })
    }

    params := &stripe.RefundParams{
        PaymentIntent: stripe.String(req.PaymentIntentID),
    }

    // 返金理由を設定（オプション）
    if req.Reason != "" {
        params.Reason = stripe.String(req.Reason)
    }

    r, err := refund.New(params)
    if err != nil {
        return c.JSON(http.StatusInternalServerError, map[string]string{
            "error": err.Error(),
        })
    }

    return c.JSON(http.StatusOK, map[string]interface{}{
        "refund": r,
        "status": r.Status,
        "amount": r.Amount,
    })
}
```

### 2.2 部分返金

```go
// 部分返金
type CreatePartialRefundRequest struct {
    PaymentIntentID string `json:"paymentIntentId"`
    Amount          int64  `json:"amount"`
    Reason          string `json:"reason,omitempty"`
}

func createPartialRefund(c echo.Context) error {
    var req CreatePartialRefundRequest
    if err := c.Bind(&req); err != nil {
        return c.JSON(http.StatusBadRequest, map[string]string{
            "error": err.Error(),
        })
    }

    params := &stripe.RefundParams{
        PaymentIntent: stripe.String(req.PaymentIntentID),
        Amount:        stripe.Int64(req.Amount),
    }

    if req.Reason != "" {
        params.Reason = stripe.String(req.Reason)
    }

    r, err := refund.New(params)
    if err != nil {
        return c.JSON(http.StatusInternalServerError, map[string]string{
            "error": err.Error(),
        })
    }

    return c.JSON(http.StatusOK, map[string]interface{}{
        "refund": r,
        "status": r.Status,
        "amount": r.Amount,
    })
}
```

### 2.3 返金履歴の取得

```go
import (
    "github.com/stripe/stripe-go/v76/charge"
)

// 特定の支払いの返金履歴を取得
func getRefundHistory(c echo.Context) error {
    chargeID := c.QueryParam("chargeId")

    params := &stripe.RefundListParams{
        Charge: stripe.String(chargeID),
    }

    i := refund.List(params)
    var refunds []*stripe.Refund

    for i.Next() {
        refunds = append(refunds, i.Refund())
    }

    if err := i.Err(); err != nil {
        return c.JSON(http.StatusInternalServerError, map[string]string{
            "error": err.Error(),
        })
    }

    return c.JSON(http.StatusOK, map[string]interface{}{
        "refunds": refunds,
    })
}
```

### 2.4 フロントエンド（返金管理画面）

```typescript
'use client';

import { useState } from 'react';

interface RefundFormProps {
  paymentIntentId: string;
  totalAmount: number;
}

export default function RefundForm({ paymentIntentId, totalAmount }: RefundFormProps) {
  const [refundType, setRefundType] = useState<'full' | 'partial'>('full');
  const [amount, setAmount] = useState(0);
  const [reason, setReason] = useState('');
  const [loading, setLoading] = useState(false);
  const [message, setMessage] = useState('');

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);

    const endpoint = refundType === 'full' 
      ? '/create-refund' 
      : '/create-partial-refund';

    const body = refundType === 'full'
      ? { paymentIntentId, reason }
      : { paymentIntentId, amount: amount * 100, reason }; // 金額をセント単位に変換

    try {
      const response = await fetch(`http://localhost:8080${endpoint}`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(body),
      });

      const data = await response.json();

      if (response.ok) {
        setMessage(`返金が完了しました。返金ID: ${data.refund.id}`);
      } else {
        setMessage(`エラー: ${data.error}`);
      }
    } catch (error) {
      setMessage('返金処理に失敗しました');
    }

    setLoading(false);
  };

  return (
    <div className="max-w-md mx-auto p-6 bg-white rounded-lg shadow">
      <h2 className="text-2xl font-bold mb-4">返金処理</h2>
      
      <form onSubmit={handleSubmit} className="space-y-4">
        <div>
          <label className="block text-sm font-medium mb-2">返金タイプ</label>
          <select
            value={refundType}
            onChange={(e) => setRefundType(e.target.value as 'full' | 'partial')}
            className="w-full border rounded px-3 py-2"
          >
            <option value="full">全額返金</option>
            <option value="partial">部分返金</option>
          </select>
        </div>

        {refundType === 'partial' && (
          <div>
            <label className="block text-sm font-medium mb-2">
              返金額（最大: ¥{totalAmount.toLocaleString()}）
            </label>
            <input
              type="number"
              value={amount}
              onChange={(e) => setAmount(Number(e.target.value))}
              max={totalAmount}
              className="w-full border rounded px-3 py-2"
              required
            />
          </div>
        )}

        <div>
          <label className="block text-sm font-medium mb-2">返金理由（オプション）</label>
          <textarea
            value={reason}
            onChange={(e) => setReason(e.target.value)}
            className="w-full border rounded px-3 py-2"
            rows={3}
          />
        </div>

        <button
          type="submit"
          disabled={loading}
          className="w-full bg-red-600 text-white py-3 rounded-lg disabled:opacity-50"
        >
          {loading ? '処理中...' : '返金する'}
        </button>

        {message && (
          <div className={`p-3 rounded ${message.includes('エラー') ? 'bg-red-100 text-red-700' : 'bg-green-100 text-green-700'}`}>
            {message}
          </div>
        )}
      </form>
    </div>
  );
}
```

---

## 請求書の自動生成

### 3.1 請求書の作成

```go
package main

import (
    "net/http"
    "time"
    "github.com/labstack/echo/v4"
    "github.com/stripe/stripe-go/v76"
    "github.com/stripe/stripe-go/v76/invoice"
    "github.com/stripe/stripe-go/v76/invoiceitem"
)

// 請求書の作成
type CreateInvoiceRequest struct {
    CustomerID  string        `json:"customerId"`
    Items       []InvoiceItem `json:"items"`
    DueDate     string        `json:"dueDate,omitempty"`
    AutoAdvance bool          `json:"autoAdvance"`
}

type InvoiceItem struct {
    Description string `json:"description"`
    Amount      int64  `json:"amount"`
    Quantity    int64  `json:"quantity"`
}

func createInvoice(c echo.Context) error {
    var req CreateInvoiceRequest
    if err := c.Bind(&req); err != nil {
        return c.JSON(http.StatusBadRequest, map[string]string{
            "error": err.Error(),
        })
    }

    // 請求書アイテムを追加
    for _, item := range req.Items {
        itemParams := &stripe.InvoiceItemParams{
            Customer:    stripe.String(req.CustomerID),
            Amount:      stripe.Int64(item.Amount),
            Currency:    stripe.String("jpy"),
            Description: stripe.String(item.Description),
            Quantity:    stripe.Int64(item.Quantity),
        }

        _, err := invoiceitem.New(itemParams)
        if err != nil {
            return c.JSON(http.StatusInternalServerError, map[string]string{
                "error": err.Error(),
            })
        }
    }

    // 請求書を作成
    invoiceParams := &stripe.InvoiceParams{
        Customer:    stripe.String(req.CustomerID),
        AutoAdvance: stripe.Bool(req.AutoAdvance),
    }

    // 支払期日を設定
    if req.DueDate != "" {
        dueDate, err := time.Parse("2006-01-02", req.DueDate)
        if err == nil {
            invoiceParams.DueDate = stripe.Int64(dueDate.Unix())
        }
    }

    inv, err := invoice.New(invoiceParams)
    if err != nil {
        return c.JSON(http.StatusInternalServerError, map[string]string{
            "error": err.Error(),
        })
    }

    return c.JSON(http.StatusOK, map[string]interface{}{
        "invoice": inv,
        "invoiceId": inv.ID,
        "hostedInvoiceUrl": inv.HostedInvoiceURL,
    })
}
```

### 3.2 請求書の確定と送信

```go
// 請求書の確定
func finalizeInvoice(c echo.Context) error {
    invoiceID := c.Param("id")

    params := &stripe.InvoiceFinalizeInvoiceParams{}
    inv, err := invoice.FinalizeInvoice(invoiceID, params)
    if err != nil {
        return c.JSON(http.StatusInternalServerError, map[string]string{
            "error": err.Error(),
        })
    }

    return c.JSON(http.StatusOK, map[string]interface{}{
        "invoice": inv,
        "status": inv.Status,
    })
}

// 請求書をメール送信
func sendInvoice(c echo.Context) error {
    invoiceID := c.Param("id")

    params := &stripe.InvoiceSendInvoiceParams{}
    inv, err := invoice.SendInvoice(invoiceID, params)
    if err != nil {
        return c.JSON(http.StatusInternalServerError, map[string]string{
            "error": err.Error(),
        })
    }

    return c.JSON(http.StatusOK, map[string]interface{}{
        "invoice": inv,
        "status": "sent",
    })
}

// 請求書の支払い
func payInvoice(c echo.Context) error {
    invoiceID := c.Param("id")

    params := &stripe.InvoicePayParams{}
    inv, err := invoice.Pay(invoiceID, params)
    if err != nil {
        return c.JSON(http.StatusInternalServerError, map[string]string{
            "error": err.Error(),
        })
    }

    return c.JSON(http.StatusOK, map[string]interface{}{
        "invoice": inv,
        "status": inv.Status,
    })
}
```

### 3.3 請求書のPDF生成とダウンロード

```go
import (
    "io"
    "net/http"
)

// 請求書PDFのダウンロード
func downloadInvoicePDF(c echo.Context) error {
    invoiceID := c.Param("id")

    // 請求書を取得
    inv, err := invoice.Get(invoiceID, nil)
    if err != nil {
        return c.JSON(http.StatusInternalServerError, map[string]string{
            "error": err.Error(),
        })
    }

    // PDFのURLを取得
    if inv.InvoicePDF == "" {
        return c.JSON(http.StatusNotFound, map[string]string{
            "error": "PDF not available",
        })
    }

    // PDFをダウンロード
    resp, err := http.Get(inv.InvoicePDF)
    if err != nil {
        return c.JSON(http.StatusInternalServerError, map[string]string{
            "error": err.Error(),
        })
    }
    defer resp.Body.Close()

    // レスポンスヘッダーを設定
    c.Response().Header().Set("Content-Type", "application/pdf")
    c.Response().Header().Set("Content-Disposition", "attachment; filename=invoice-"+invoiceID+".pdf")

    // PDFをストリーム
    _, err = io.Copy(c.Response().Writer, resp.Body)
    if err != nil {
        return err
    }

    return nil
}
```

### 3.4 フロントエンド（請求書作成フォーム）

```typescript
'use client';

import { useState } from 'react';

interface InvoiceItem {
  description: string;
  amount: number;
  quantity: number;
}

export default function CreateInvoiceForm({ customerId }: { customerId: string }) {
  const [items, setItems] = useState<InvoiceItem[]>([
    { description: '', amount: 0, quantity: 1 },
  ]);
  const [dueDate, setDueDate] = useState('');
  const [autoAdvance, setAutoAdvance] = useState(true);
  const [loading, setLoading] = useState(false);
  const [result, setResult] = useState<any>(null);

  const addItem = () => {
    setItems([...items, { description: '', amount: 0, quantity: 1 }]);
  };

  const updateItem = (index: number, field: keyof InvoiceItem, value: any) => {
    const newItems = [...items];
    newItems[index] = { ...newItems[index], [field]: value };
    setItems(newItems);
  };

  const removeItem = (index: number) => {
    setItems(items.filter((_, i) => i !== index));
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);

    try {
      const response = await fetch('http://localhost:8080/create-invoice', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          customerId,
          items: items.map(item => ({
            ...item,
            amount: item.amount * 100, // セント単位に変換
          })),
          dueDate,
          autoAdvance,
        }),
      });

      const data = await response.json();
      setResult(data);
    } catch (error) {
      console.error('Error:', error);
    }

    setLoading(false);
  };

  const totalAmount = items.reduce(
    (sum, item) => sum + item.amount * item.quantity,
    0
  );

  return (
    <div className="max-w-3xl mx-auto p-6">
      <h2 className="text-2xl font-bold mb-6">請求書作成</h2>

      <form onSubmit={handleSubmit} className="space-y-6">
        {/* 請求書アイテム */}
        <div>
          <h3 className="text-lg font-semibold mb-3">請求項目</h3>
          {items.map((item, index) => (
            <div key={index} className="grid grid-cols-12 gap-3 mb-3">
              <input
                type="text"
                placeholder="項目名"
                value={item.description}
                onChange={(e) => updateItem(index, 'description', e.target.value)}
                className="col-span-5 border rounded px-3 py-2"
                required
              />
              <input
                type="number"
                placeholder="単価"
                value={item.amount || ''}
                onChange={(e) => updateItem(index, 'amount', Number(e.target.value))}
                className="col-span-3 border rounded px-3 py-2"
                required
              />
              <input
                type="number"
                placeholder="数量"
                value={item.quantity || ''}
                onChange={(e) => updateItem(index, 'quantity', Number(e.target.value))}
                className="col-span-2 border rounded px-3 py-2"
                required
                min="1"
              />
              <button
                type="button"
                onClick={() => removeItem(index)}
                className="col-span-2 bg-red-500 text-white rounded px-3 py-2"
                disabled={items.length === 1}
              >
                削除
              </button>
            </div>
          ))}
          <button
            type="button"
            onClick={addItem}
            className="bg-gray-200 px-4 py-2 rounded"
          >
            + 項目を追加
          </button>
        </div>

        {/* 合計金額 */}
        <div className="bg-gray-100 p-4 rounded">
          <p className="text-xl font-bold">
            合計: ¥{totalAmount.toLocaleString()}
          </p>
        </div>

        {/* 支払期日 */}
        <div>
          <label className="block text-sm font-medium mb-2">
            支払期日（オプション）
          </label>
          <input
            type="date"
            value={dueDate}
            onChange={(e) => setDueDate(e.target.value)}
            className="border rounded px-3 py-2"
          />
        </div>

        {/* 自動確定 */}
        <div className="flex items-center">
          <input
            type="checkbox"
            checked={autoAdvance}
            onChange={(e) => setAutoAdvance(e.target.checked)}
            className="mr-2"
          />
          <label className="text-sm">
            自動的に請求書を確定して送信する
          </label>
        </div>

        <button
          type="submit"
          disabled={loading}
          className="w-full bg-blue-600 text-white py-3 rounded-lg disabled:opacity-50"
        >
          {loading ? '作成中...' : '請求書を作成'}
        </button>
      </form>

      {/* 結果表示 */}
      {result && (
        <div className="mt-6 p-4 bg-green-100 rounded">
          <p className="font-semibold">請求書が作成されました！</p>
          <p className="text-sm mt-2">請求書ID: {result.invoiceId}</p>
          {result.hostedInvoiceUrl && (
            <a
              href={result.hostedInvoiceUrl}
              target="_blank"
              rel="noopener noreferrer"
              className="text-blue-600 underline mt-2 inline-block"
            >
              請求書を表示
            </a>
          )}
        </div>
      )}
    </div>
  );
}
```

### 3.5 請求書一覧の表示

```typescript
'use client';

import { useState, useEffect } from 'react';

interface Invoice {
  id: string;
  number: string;
  status: string;
  amount_due: number;
  due_date: number;
  created: number;
  hosted_invoice_url: string;
  invoice_pdf: string;
}

export default function InvoiceList({ customerId }: { customerId: string }) {
  const [invoices, setInvoices] = useState<Invoice[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    fetch(`http://localhost:8080/invoices?customerId=${customerId}`)
      .then((res) => res.json())
      .then((data) => {
        setInvoices(data.invoices);
        setLoading(false);
      });
  }, [customerId]);

  const handleDownloadPDF = async (invoiceId: string) => {
    const response = await fetch(
      `http://localhost:8080/invoice/${invoiceId}/pdf`
    );
    const blob = await response.blob();
    const url = window.URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url;
    a.download = `invoice-${invoiceId}.pdf`;
    a.click();
  };

  const getStatusBadge = (status: string) => {
    const styles = {
      paid: 'bg-green-100 text-green-800',
      open: 'bg-yellow-100 text-yellow-800',
      draft: 'bg-gray-100 text-gray-800',
      void: 'bg-red-100 text-red-800',
    };
    return styles[status as keyof typeof styles] || 'bg-gray-100 text-gray-800';
  };

  if (loading) return <div>読み込み中...</div>;

  return (
    <div className="max-w-4xl mx-auto p-6">
      <h2 className="text-2xl font-bold mb-6">請求書一覧</h2>
      
      <div className="space-y-4">
        {invoices.map((invoice) => (
          <div
            key={invoice.id}
            className="border rounded-lg p-4 flex justify-between items-center"
          >
            <div className="flex-1">
              <div className="flex items-center gap-3">
                <p className="font-semibold">#{invoice.number}</p>
                <span
                  className={`px-2 py-1 rounded text-xs ${getStatusBadge(
                    invoice.status
                  )}`}
                >
                  {invoice.status}
                </span>
              </div>
              <p className="text-sm text-gray-600 mt-1">
                金額: ¥{(invoice.amount_due / 100).toLocaleString()}
              </p>
              {invoice.due_date && (
                <p className="text-sm text-gray-600">
                  期日: {new Date(invoice.due_date * 1000).toLocaleDateString()}
                </p>
              )}
            </div>
            
            <div className="flex gap-2">
              {invoice.hosted_invoice_url && (
                <a
                  href={invoice.hosted_invoice_url}
                  target="_blank"
                  rel="noopener noreferrer"
                  className="bg-blue-600 text-white px-4 py-2 rounded"
                >
                  表示
                </a>
              )}
              {invoice.invoice_pdf && (
                <button
                  onClick={() => handleDownloadPDF(invoice.id)}
                  className="bg-gray-600 text-white px-4 py-2 rounded"
                >
                  PDF
                </button>
              )}
            </div>
          </div>
        ))}
      </div>
    </div>
  );
}
```

---

## 実践的な実装例

### 4.1 統合的なルーティング設定

```go
func main() {
    if err := godotenv.Load(); err != nil {
        log.Fatal("Error loading .env file")
    }

    stripe.Key = os.Getenv("STRIPE_SECRET_KEY")

    e := echo.New()

    e.Use(middleware.Logger())
    e.Use(middleware.Recover())
    e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
        AllowOrigins: []string{"http://localhost:3000"},
        AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodPUT, http.MethodDELETE},
    }))

    // 顧客管理
    e.POST("/customers", createCustomer)
    
    // 支払い
    e.POST("/create-payment-intent", createPaymentIntent)
    e.POST("/create-checkout-session", createCheckoutSession)
    
    // 支払い方法管理
    e.POST("/attach-payment-method", attachPaymentMethod)
    e.GET("/payment-methods", listPaymentMethods)
    e.POST("/set-default-payment-method", setDefaultPaymentMethod)
    
    // 返金
    e.POST("/create-refund", createRefund)
    e.POST("/create-partial-refund", createPartialRefund)
    e.GET("/refunds", getRefundHistory)
    
    // 請求書
    e.POST("/create-invoice", createInvoice)
    e.POST("/invoice/:id/finalize", finalizeInvoice)
    e.POST("/invoice/:id/send", sendInvoice)
    e.POST("/invoice/:id/pay", payInvoice)
    e.GET("/invoice/:id/pdf", downloadInvoicePDF)
    e.GET("/invoices", listInvoices)
    
    // Webhook
    e.POST("/webhook", handleWebhook)

    e.Logger.Fatal(e.Start(":8080"))
}
```

### 4.2 Webhookイベントハンドリングの拡張

```go
func handleWebhook(c echo.Context) error {
    payload, err := io.ReadAll(c.Request().Body)
    if err != nil {
        return c.JSON(http.StatusBadRequest, map[string]string{
            "error": "Invalid payload",
        })
    }

    signature := c.Request().Header.Get("Stripe-Signature")
    webhookSecret := os.Getenv("STRIPE_WEBHOOK_SECRET")

    event, err := webhook.ConstructEvent(payload, signature, webhookSecret)
    if err != nil {
        log.Printf("Webhook signature verification failed: %v", err)
        return c.JSON(http.StatusBadRequest, map[string]string{
            "error": err.Error(),
        })
    }

    switch event.Type {
    case "payment_intent.succeeded":
        log.Println("Payment intent succeeded")
        // 支払い成功時の処理
        
    case "payment_intent.payment_failed":
        log.Println("Payment intent failed")
        // 支払い失敗時の処理
        
    case "charge.refunded":
        log.Println("Charge refunded")
        // 返金時の処理
        
    case "invoice.created":
        log.Println("Invoice created")
        // 請求書作成時の処理
        
    case "invoice.finalized":
        log.Println("Invoice finalized")
        // 請求書確定時の処理
        
    case "invoice.paid":
        log.Println("Invoice paid")
        // 請求書支払い時の処理
        
    case "invoice.payment_failed":
        log.Println("Invoice payment failed")
        // 請求書支払い失敗時の処理
        
    case "customer.created":
        log.Println("Customer created")
        // 顧客作成時の処理
        
    case "payment_method.attached":
        log.Println("Payment method attached")
        // 支払い方法追加時の処理
        
    default:
        log.Printf("Unhandled event type: %s", event.Type)
    }

    return c.JSON(http.StatusOK, map[string]string{
        "status": "success",
    })
}
```

---

## まとめ

このドキュメントでは、Stripeの応用機能として以下を実装しました。

### 実装した機能

1. **複数の支払い方法のサポート**
   - Payment Intentを使用した柔軟な支払い
   - 保存済み支払い方法の管理
   - デフォルト支払い方法の設定

2. **返金処理**
   - 全額返金・部分返金
   - 返金履歴の管理
   - 管理画面での返金操作

3. **請求書の自動生成**
   - 請求書の作成と管理
   - PDF生成とダウンロード
   - メール送信機能

### セキュリティのベストプラクティス

- Webhook署名の検証を必ず実施
- APIキーを環境変数で管理
- HTTPS通信の使用
- エラーハンドリングの適切な実装
- ログの記録と監視

### 参考リンク

- [Stripe Payment Intents API](https://stripe.com/docs/api/payment_intents)
- [Stripe Refunds](https://stripe.com/docs/refunds)
- [Stripe Invoicing](https://stripe.com/docs/invoicing)
- [Stripe Webhooks](https://stripe.com/docs/webhooks)
