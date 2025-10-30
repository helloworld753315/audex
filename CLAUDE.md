# audex (Audio Extractor)

MP4ファイルから音声のみを抽出してM4Aファイルに変換するGoツール。
アルバムアートワークとメタデータを完全に保持します。

**プロジェクト名:** audex (Audio Extractor)
**読み方:** オーデックス

## プロジェクト概要

### 目的
MP4ファイル（音声のみ）を、適切な音声ファイル形式であるM4Aに変換する。
ffmpeg コマンドをシェルスクリプトで叩いていたものをGoに移行したいです。
https://github.com/helloworld753315/ffmpeg-scripts

### 主な特徴
- **ffmpegに依存しない** - 純粋なGo実装
- **無劣化変換** - ストリームコピー方式（再エンコードなし）
- **超高速** - 音声データをそのままコピーするだけ
- **メタデータ完全保持** - タイトル、アーティスト、アルバム等
- **アルバムアートワーク保持** - ジャケット画像を完全保持
- **MITライセンス** - オープンソースとして公開予定

## 技術スタック

### 言語
- Go 1.21以上

### 主要ライブラリ

**Phase 1 (必須):**
- `github.com/abema/go-mp4` (MIT) - MP4コンテナ操作、音声ストリーム抽出
- `github.com/dhowden/tag` (Unlicense) - メタデータ・アートワーク読み書き
- `github.com/urfave/cli/v2` (MIT) - CLIフレームワーク

**Phase 2 (機能拡張時):**
- `github.com/schollz/progressbar/v3` (MIT) - プログレスバー表示
- シェル補完スクリプト生成機能（urfave/cli組み込み）

## アーキテクチャ

### 処理フロー
```
1. 入力MP4ファイルを開く
2. tagライブラリでメタデータを読み取る
   - タイトル、アーティスト、アルバム、年
   - アルバムアートワーク（JPEG/PNG）
3. go-mp4で音声トラックのみを抽出
   - ftypボックスを"M4A "に設定
   - 映像トラックを除外
   - 音声トラックをストリームコピー
   - メタデータボックスを保持
4. M4Aファイルとして出力
5. tagライブラリでメタデータを書き込む
```

### 重要な技術的決定

#### なぜストリームコピーか
- 再エンコード不要 → 音質劣化なし
- 処理が超高速（数秒）
- CPU負荷が低い
- 元のAAC品質を完全保持

#### なぜgo-mp4 + tagの組み合わせか
- go-mp4: MP4コンテナ操作が得意
- tag: メタデータ処理が簡単
- 両方とも純粋なGo実装
- クロスコンパイルが簡単
- 依存関係がクリーン

## MVP（最小機能製品）の要件

### Phase 1: 基本機能
- [ ] 単一ファイル変換（MP4 → M4A）
- [ ] 音声ストリームコピー（go-mp4）
- [ ] メタデータ保持（tag）
  - [ ] タイトル
  - [ ] アーティスト
  - [ ] アルバム
  - [ ] 年
  - [ ] トラック番号
  - [ ] ジャンル
- [ ] アルバムアートワーク保持
- [ ] 基本的なエラーハンドリング
- [ ] CLIインターフェース: `audex input.mp4 output.m4a`

### Phase 2: 機能拡張（後回し）
- [ ] バッチ処理（複数ファイル一括変換）
- [ ] 並列処理（goroutine）
- [ ] プログレスバー
- [ ] カラー出力
- [ ] ディレクトリ再帰処理
- [ ] ドライラン（--dry-run）
- [ ] 詳細なログ出力（--verbose）

### Phase 3: 将来的な拡張（検討中）
- [ ] MP3出力対応（再エンコード必要、外部依存あり）
- [ ] FLAC対応
- [ ] GUI版

## 技術的な課題と解決策

### 課題1: go-mp4でのメタデータ処理が複雑
**解決策**: tagライブラリを併用
- go-mp4は音声抽出のみに使用
- メタデータはtagで処理
- 責任を明確に分離

### 課題2: mdatボックスの処理
**問題**: mdatには映像と音声が混在
**解決策**: 
- サンプルテーブル(stbl)を参照
- 音声サンプルの位置を特定
- 音声データのみをコピー

### 課題3: Appleメタデータフォーマット
**解決策**: tagライブラリがApple iTunesメタデータ（ilst）に対応済み

## ファイル構造

```
audex/
├── CLAUDE.md           # このファイル
├── README.md           # ユーザー向けドキュメント
├── LICENSE             # MITライセンス
├── go.mod
├── go.sum
├── main.go             # エントリーポイント
├── converter/
│   ├── converter.go    # 変換ロジック
│   ├── metadata.go     # メタデータ処理
│   └── stream.go       # ストリーム操作
├── internal/
│   └── mp4/
│       └── extractor.go # MP4音声抽出
└── cmd/
    └── audex/
        └── main.go
```

## 開発ガイドライン

### コーディング規約
- 標準的なGo規約に従う
- `gofmt`でフォーマット
- エラーは明確に返す
- ログは標準出力/エラー出力を使用

### テスト戦略

#### テストの3層構造（段階的に実装）

**1. ユニットテスト（Unit Tests）** ← **Phase 1: まずはここだけ**
```
目的: 個別関数の正確性を保証
範囲: 各パッケージの関数単位
実行: 高速（秒単位）
優先度: ★★★ 必須
```

**2. 統合テスト（Integration Tests）** ← **Phase 3: 余裕があれば**
```
目的: ライブラリ連携の動作確認
範囲: go-mp4とtagの組み合わせ
実行: 中速（数秒〜数十秒）
優先度: ★★☆ 重要（後回しOK）
```

**3. E2Eテスト（End-to-End Tests）** ← **Phase 4: リリース前**
```
目的: 実際のユースケースを検証
範囲: 実ファイルでの変換全体
実行: 低速（数十秒〜数分）
優先度: ★☆☆ あれば良い
```

**Phase 2（手動テスト）も重要:**
```
実際のSONYリッピングMP4で試す
→ 自動テストより先に動作確認
→ 問題があればすぐに気づける
```

#### テストデータ

```
testdata/
├── valid/                    # 正常系テストデータ
│   ├── simple.mp4           # 最小限のMP4（音声のみ）
│   ├── with_artwork.mp4     # アートワーク付き
│   ├── full_metadata.mp4    # 全メタデータ付き
│   ├── sony_ripped.mp4      # SONYリッピング形式
│   └── expected/            # 期待される出力
│       ├── simple.m4a
│       └── ...
├── invalid/                  # 異常系テストデータ
│   ├── corrupted.mp4        # 壊れたMP4
│   ├── video_only.mp4       # 映像のみ（音声なし）
│   ├── no_audio.mp4         # 音声トラックなし
│   └── encrypted.mp4        # 暗号化ファイル
└── edge_cases/              # エッジケース
    ├── large_file.mp4       # 大容量ファイル（500MB+）
    ├── long_title.mp4       # 超長いメタデータ
    ├── special_chars.mp4    # 特殊文字を含む
    └── multi_audio.mp4      # 複数音声トラック
```

**テストデータの作成方法:**
```bash
# ffmpegで最小限のテストMP4を作成（開発初期のみ使用）
ffmpeg -f lavfi -i sine=frequency=440:duration=1 -c:a aac testdata/valid/simple.mp4

# 後はツール自体でテストデータを作成
# 本番環境からサンプルファイルを収集（個人情報削除）
```

#### ユニットテストの設計

**converter/metadata_test.go**
```go
package converter

import (
	"testing"
)

func TestReadMetadata(t *testing.T) {
	tests := []struct {
		name     string
		file     string
		want     *Metadata
		wantErr  bool
	}{
		{
			name: "有効なメタデータ",
			file: "../testdata/valid/full_metadata.mp4",
			want: &Metadata{
				Title:  "Test Song",
				Artist: "Test Artist",
				Album:  "Test Album",
				Year:   "2024",
			},
			wantErr: false,
		},
		{
			name:    "メタデータなし",
			file:    "../testdata/valid/simple.mp4",
			want:    &Metadata{},
			wantErr: false,
		},
		{
			name:    "ファイルが存在しない",
			file:    "nonexistent.mp4",
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ReadMetadata(tt.file)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadMetadata() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !equalMetadata(got, tt.want) {
				t.Errorf("ReadMetadata() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReadArtwork(t *testing.T) {
	tests := []struct {
		name     string
		file     string
		wantSize int  // アートワークのサイズ（バイト）
		wantMIME string
		wantErr  bool
	}{
		{
			name:     "JPEGアートワーク",
			file:     "../testdata/valid/with_artwork.mp4",
			wantSize: 10000, // 約10KB
			wantMIME: "image/jpeg",
			wantErr:  false,
		},
		{
			name:     "アートワークなし",
			file:     "../testdata/valid/simple.mp4",
			wantSize: 0,
			wantMIME: "",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			artwork, err := ReadArtwork(tt.file)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadArtwork() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && artwork != nil {
				if len(artwork.Data) < tt.wantSize/2 {
					t.Errorf("Artwork too small: got %d bytes, want at least %d", 
						len(artwork.Data), tt.wantSize/2)
				}
				if artwork.MIMEType != tt.wantMIME {
					t.Errorf("Wrong MIME type: got %s, want %s", 
						artwork.MIMEType, tt.wantMIME)
				}
			}
		})
	}
}
```

**converter/stream_test.go**
```go
package converter

import (
	"os"
	"testing"
)

func TestExtractAudioStream(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{
			name:    "正常な音声抽出",
			input:   "../testdata/valid/simple.mp4",
			wantErr: false,
		},
		{
			name:    "音声トラックなし",
			input:   "../testdata/invalid/video_only.mp4",
			wantErr: true,
		},
		{
			name:    "壊れたファイル",
			input:   "../testdata/invalid/corrupted.mp4",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := t.TempDir() + "/output.m4a"
			err := ExtractAudioStream(tt.input, output)
			if (err != nil) != tt.wantErr {
				t.Errorf("ExtractAudioStream() error = %v, wantErr %v", err, tt.wantErr)
			}
			
			// 正常系の場合、ファイルが作成されているか確認
			if !tt.wantErr {
				if _, err := os.Stat(output); os.IsNotExist(err) {
					t.Error("Output file was not created")
				}
			}
		})
	}
}
```

#### 統合テストの設計

**integration_test.go**
```go
package main

import (
	"os"
	"testing"
	
	"github.com/dhowden/tag"
)

func TestFullConversion(t *testing.T) {
	// 統合テスト: メタデータ保持を含む完全な変換
	
	input := "testdata/valid/full_metadata.mp4"
	output := t.TempDir() + "/output.m4a"
	
	// 変換前のメタデータを記録
	originalMetadata := readMetadataForTest(t, input)
	
	// 変換実行
	err := convert(input, output)
	if err != nil {
		t.Fatalf("Conversion failed: %v", err)
	}
	
	// 変換後のメタデータを確認
	convertedMetadata := readMetadataForTest(t, output)
	
	// メタデータが保持されているか検証
	if originalMetadata.Title != convertedMetadata.Title {
		t.Errorf("Title not preserved: got %s, want %s", 
			convertedMetadata.Title, originalMetadata.Title)
	}
	
	if originalMetadata.Artist != convertedMetadata.Artist {
		t.Errorf("Artist not preserved: got %s, want %s", 
			convertedMetadata.Artist, originalMetadata.Artist)
	}
	
	// アートワークが保持されているか検証
	if originalMetadata.Picture() != nil && convertedMetadata.Picture() == nil {
		t.Error("Artwork was not preserved")
	}
}

func TestArtworkPreservation(t *testing.T) {
	// アートワーク保持の統合テスト
	
	input := "testdata/valid/with_artwork.mp4"
	output := t.TempDir() + "/output.m4a"
	
	// 元のアートワーク情報
	f, _ := os.Open(input)
	defer f.Close()
	m, _ := tag.ReadFrom(f)
	originalPic := m.Picture()
	
	// 変換
	if err := convert(input, output); err != nil {
		t.Fatalf("Conversion failed: %v", err)
	}
	
	// 変換後のアートワーク確認
	f2, _ := os.Open(output)
	defer f2.Close()
	m2, _ := tag.ReadFrom(f2)
	convertedPic := m2.Picture()
	
	if convertedPic == nil {
		t.Fatal("Artwork was lost during conversion")
	}
	
	// サイズが同じか確認
	if len(convertedPic.Data) != len(originalPic.Data) {
		t.Errorf("Artwork size mismatch: got %d, want %d", 
			len(convertedPic.Data), len(originalPic.Data))
	}
	
	// MIMEタイプが同じか確認
	if convertedPic.MIMEType != originalPic.MIMEType {
		t.Errorf("Artwork MIME type mismatch: got %s, want %s", 
			convertedPic.MIMEType, originalPic.MIMEType)
	}
}

func readMetadataForTest(t *testing.T, path string) tag.Metadata {
	f, err := os.Open(path)
	if err != nil {
		t.Fatalf("Failed to open file: %v", err)
	}
	defer f.Close()
	
	m, err := tag.ReadFrom(f)
	if err != nil {
		t.Fatalf("Failed to read metadata: %v", err)
	}
	
	return m
}
```

#### E2Eテストの設計

**e2e_test.go**
```go
//go:build e2e
// +build e2e

package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

// E2Eテストは実際のバイナリを使用
// go test -tags=e2e ./...

func TestE2E_BasicConversion(t *testing.T) {
	// バイナリをビルド
	binary := buildBinary(t)
	defer os.Remove(binary)
	
	// テスト実行
	input := "testdata/valid/sony_ripped.mp4"
	output := t.TempDir() + "/output.m4a"
	
	cmd := exec.Command(binary, input, output)
	output_bytes, err := cmd.CombinedOutput()
	
	if err != nil {
		t.Fatalf("Command failed: %v\nOutput: %s", err, output_bytes)
	}
	
	// 出力ファイルの存在確認
	if _, err := os.Stat(output); err != nil {
		t.Errorf("Output file not created: %v", err)
	}
	
	// ファイルサイズの確認（元ファイルより小さいはず）
	inputInfo, _ := os.Stat(input)
	outputInfo, _ := os.Stat(output)
	
	if outputInfo.Size() > inputInfo.Size() {
		t.Errorf("Output file larger than input: %d > %d", 
			outputInfo.Size(), inputInfo.Size())
	}
}

func TestE2E_BatchConversion(t *testing.T) {
	// バッチ処理のE2Eテスト（Phase 2実装後）
	t.Skip("Not implemented yet")
}

func TestE2E_ErrorHandling(t *testing.T) {
	binary := buildBinary(t)
	defer os.Remove(binary)
	
	tests := []struct {
		name    string
		args    []string
		wantErr bool
	}{
		{
			name:    "引数なし",
			args:    []string{},
			wantErr: true,
		},
		{
			name:    "存在しないファイル",
			args:    []string{"nonexistent.mp4", "output.m4a"},
			wantErr: true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := exec.Command(binary, tt.args...)
			err := cmd.Run()
			
			if tt.wantErr && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.wantErr && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}

func buildBinary(t *testing.T) string {
	binary := filepath.Join(t.TempDir(), "mp4tom4a")
	cmd := exec.Command("go", "build", "-o", binary)
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to build binary: %v", err)
	}
	return binary
}
```

#### ベンチマークテスト

**benchmark_test.go**
```go
package converter

import (
	"testing"
)

func BenchmarkExtractAudioStream(b *testing.B) {
	input := "../testdata/valid/simple.mp4"
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		output := b.TempDir() + "/output.m4a"
		ExtractAudioStream(input, output)
	}
}

func BenchmarkMetadataRead(b *testing.B) {
	input := "../testdata/valid/full_metadata.mp4"
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ReadMetadata(input)
	}
}

func BenchmarkFullConversion(b *testing.B) {
	input := "../testdata/valid/sony_ripped.mp4"
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		output := b.TempDir() + "/output.m4a"
		convert(input, output)
	}
}
```

#### テスト実行コマンド

```bash
# すべてのユニットテスト
go test ./...

# カバレッジ付き
go test -cover ./...

# 詳細なカバレッジレポート
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# 統合テストも含める
go test -v ./...

# E2Eテストのみ
go test -tags=e2e ./...

# ベンチマーク
go test -bench=. ./...

# 並列実行
go test -parallel 4 ./...
```

#### CI/CDでのテスト戦略

```yaml
# .github/workflows/test.yml
name: Test

on: [push, pull_request]

jobs:
  test:
    strategy:
      matrix:
        os: [ubuntu-latest, macos-latest, windows-latest]
        go: ['1.21', '1.22']
    
    runs-on: ${{ matrix.os }}
    
    steps:
      - uses: actions/checkout@v3
      
      - name: Set up Go
        uses: actions/setup-go@v4
        with:
          go-version: ${{ matrix.go }}
      
      - name: Unit Tests
        run: go test -v ./...
      
      - name: Coverage
        run: go test -coverprofile=coverage.out ./...
      
      - name: Upload Coverage
        uses: codecov/codecov-action@v3
        with:
          file: ./coverage.out
      
      - name: E2E Tests
        run: go test -tags=e2e -v ./...
```

#### テストカバレッジ目標

```
Phase 1 (MVP - ユニットテストのみ):
- ユニットテスト: 40-60% でOK
- 重要な関数だけカバー
- 統合・E2Eテスト: なし（手動テストで代用）

Phase 2 (機能追加後):
- ユニットテスト: 80%以上
- 統合テスト: 主要機能をカバー
- E2Eテスト: 基本シナリオのみ

長期目標:
- ユニットテスト: 90%以上
- エッジケース完全カバー
```

#### MVP（Phase 1）での現実的なテスト方針

**書くべきユニットテスト:**
```
✅ メタデータ読み取り（tag使用）
✅ メタデータ書き込み（tag使用）
✅ エラーハンドリング（ファイル存在チェックなど）
✅ シンプルなヘルパー関数
```

**書かなくてOK（Phase 1では）:**
```
⏭️ go-mp4の複雑な部分（ライブラリ自体がテスト済み）
⏭️ 統合テスト（手動でやる）
⏭️ E2Eテスト（リリース前に追加）
⏭️ ベンチマーク（最適化時に追加）
```

**手動テストで確認:**
```
1. 実際のSONYリッピングMP4で変換
2. メタデータが保持されているか目視確認
3. アートワークが表示されるか確認
4. 音声が正常に再生されるか確認
```

**この方針の理由:**
- まず動くものを作ることが最優先
- テストに時間をかけすぎない
- ユーザーフィードバックを早く得る
- 必要に応じてテストを追加

#### テストの優先順位

**高優先度（必須）:**
1. ✅ メタデータ読み書き
2. ✅ アートワーク保持
3. ✅ 音声ストリーム抽出
4. ✅ エラーハンドリング

**中優先度（重要）:**
5. ⏺️ 特殊文字対応
6. ⏺️ 大容量ファイル対応
7. ⏺️ 並列処理

**低優先度（あれば良い）:**
8. ⏺️ パフォーマンス最適化
9. ⏺️ メモリ使用量最適化

### ドキュメント
- 各関数にGoDocコメント
- 複雑なロジックにはインラインコメント
- READMEに使用例を記載

## 実装の優先順位

### 1. まず動くものを作る（1-2日）
```go
// 最小限の実装
func main() {
    input := os.Args[1]
    output := os.Args[2]
    
    // 1. メタデータ読み取り（tag）
    metadata := readMetadata(input)
    
    // 2. 音声ストリーム抽出（go-mp4）
    extractAudio(input, output)
    
    // 3. メタデータ書き込み（tag）
    writeMetadata(output, metadata)
}
```

### 2. エラーハンドリング追加（1日）
- ファイル存在チェック
- フォーマット検証
- 詳細なエラーメッセージ

### 3. リファクタリング（1日）
- コードを整理
- パッケージ分割
- テスト追加

### 4. 機能拡張（必要に応じて）
- バッチ処理
- 並列処理
- CLI改善

## 参考情報

### MP4/M4Aフォーマット
- ISO/IEC 14496-12 (ISO Base Media File Format)
- ISO/IEC 14496-14 (MP4 File Format)
- Apple iTunes Metadata Format

### 既存の実装（参考）
- ffmpeg（C実装、GPL/LGPL）
- 現在のシェルスクリプト版（リポジトリ: helloworld753315/ffmpeg-scripts）

## ライセンス

MITライセンスで公開予定

## 開発メモ

### なぜこのツールを作るのか
- ffmpegは高機能だが複雑
- 単一バイナリで配布したい
- 特定用途（CDリッピングMP4→M4A）に最適化
- Go言語の学習も兼ねて

### ターゲットユーザー
- SONYのCDリッピングソフトユーザー
- MP4で音声を保存している人
- 音楽ファイルを整理したい人
- ffmpegをインストールしたくない人

## Claude Codeへの指示

このプロジェクトを実装する際の注意点：

1. **段階的に実装してください**
   - まず最小限の動くものを作成
   - 機能を一つずつ追加
   - 各段階で動作確認

2. **go-mp4とtagの使い方を調べてから実装**
   - ライブラリのドキュメントを確認
   - サンプルコードを参考に
   - 両方のライブラリの特性を理解

3. **エラーハンドリングを重視**
   - ユーザーフレンドリーなエラーメッセージ
   - 適切なエラーチェック
   - パニックを避ける

4. **テスト可能な設計**
   - 関数を小さく保つ
   - 依存を注入可能に
   - モックしやすい構造

5. **コメントを適切に**
   - なぜそうしたのかを説明
   - 複雑な処理には詳細な説明
   - TODO/FIXMEを明確に

## 最初のタスク

Claude Codeで実装を始める際は、以下の順序で進めてください：

1. ✅ プロジェクト構造を作成（go.mod, main.go）
2. ⬜ tagライブラリでメタデータ読み取り機能を実装
3. ⬜ tagライブラリでメタデータ書き込み機能を実装
4. ⬜ go-mp4で音声ストリーム抽出機能を実装（最難関）
5. ⬜ 全体を統合して動作確認
6. ⬜ エラーハンドリング追加
7. ⬜ リファクタリングとテスト

まずは2から始めましょう。tagライブラリは比較的簡単なので、
メタデータの読み書きを実装して、その後go-mp4の複雑な部分に取り組むのが良いでしょう。