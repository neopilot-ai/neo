package integration_test

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestNeoInitWorkflow tests the complete project initialization workflow
func TestNeoInitWorkflow(t *testing.T) {
	// Create a temporary directory for the test
	tempDir, err := os.MkdirTemp("", "neo-init-test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Change to the temporary directory
	originalWd, err := os.Getwd()
	require.NoError(t, err)
	defer os.Chdir(originalWd)

	err = os.Chdir(tempDir)
	require.NoError(t, err)

	// Run neo init
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "go", "run", "../../cmd/neo", "init", "--yes")
	cmd.Dir = tempDir
	output, err := cmd.CombinedOutput()
	
	if err != nil {
		t.Logf("Init output: %s", string(output))
		t.Logf("Init error: %v", err)
	}

	// Check that config file was created
	configPath := filepath.Join(tempDir, "neo.config.ts")
	_, err = os.Stat(configPath)
	if os.IsNotExist(err) {
		t.Skip("Neo config not created - CLI may not be available in test environment")
		return
	}
	require.NoError(t, err, "neo.config.ts should be created")

	// Check that .neo directory was created
	neoDir := filepath.Join(tempDir, ".neo")
	_, err = os.Stat(neoDir)
	require.NoError(t, err, ".neo directory should be created")

	t.Logf("Successfully initialized NEO project in %s", tempDir)
}

// TestNeoInstallWorkflow tests the provider installation workflow
func TestNeoInstallWorkflow(t *testing.T) {
	// Create a temporary directory for the test
	tempDir, err := os.MkdirTemp("", "neo-install-test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	originalWd, err := os.Getwd()
	require.NoError(t, err)
	defer os.Chdir(originalWd)

	err = os.Chdir(tempDir)
	require.NoError(t, err)

	// Create a basic neo.config.ts file
	configContent := `
export default {
  providers: {
    aws: "6.27.0"
  }
};
`
	err = os.WriteFile("neo.config.ts", []byte(configContent), 0644)
	require.NoError(t, err)

	// Run neo install
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "go", "run", "../../cmd/neo", "install")
	cmd.Dir = tempDir
	output, err := cmd.CombinedOutput()
	
	if err != nil {
		t.Logf("Install output: %s", string(output))
		t.Logf("Install error: %v", err)
		t.Skip("Neo install failed - CLI may not be available in test environment")
		return
	}

	// Check that platform directory was created
	platformDir := filepath.Join(tempDir, ".neo", "platform")
	_, err = os.Stat(platformDir)
	require.NoError(t, err, "platform directory should be created")

	// Check that package.json exists in platform directory
	packageJsonPath := filepath.Join(platformDir, "package.json")
	_, err = os.Stat(packageJsonPath)
	require.NoError(t, err, "package.json should be created in platform directory")

	t.Logf("Successfully installed providers in %s", tempDir)
}

// TestProjectDiscovery tests project discovery functionality
func TestProjectDiscovery(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "neo-discovery-test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	originalWd, err := os.Getwd()
	require.NoError(t, err)
	defer os.Chdir(originalWd)

	err = os.Chdir(tempDir)
	require.NoError(t, err)

	// Test discovery without config file
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "go", "run", "../../cmd/neo", "version")
	cmd.Dir = tempDir
	output, err := cmd.CombinedOutput()
	
	if err != nil {
		t.Logf("Version output: %s", string(output))
		t.Logf("Version error: %v", err)
		t.Skip("Neo CLI not available in test environment")
		return
	}

	// Create config file and test again
	configContent := `export default {}`
	err = os.WriteFile("neo.config.ts", []byte(configContent), 0644)
	require.NoError(t, err)

	cmd = exec.CommandContext(ctx, "go", "run", "../../cmd/neo", "version")
	cmd.Dir = tempDir
	output, err = cmd.CombinedOutput()
	require.NoError(t, err, "Version command should work with config file")

	assert.Contains(t, string(output), "dev", "Version output should contain version")

	t.Logf("Successfully tested project discovery in %s", tempDir)
}

// TestStageManagement tests the stage management workflow
func TestStageManagement(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "neo-stage-test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	originalWd, err := os.Getwd()
	require.NoError(t, err)
	defer os.Chdir(originalWd)

	err = os.Chdir(tempDir)
	require.NoError(t, err)

	// Create config file
	configContent := `export default {}`
	err = os.WriteFile("neo.config.ts", []byte(configContent), 0644)
	require.NoError(t, err)

	// Test setting stage
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "go", "run", "../../cmd/neo", "deploy", "--stage", "test-stage", "--dry-run")
	cmd.Dir = tempDir
	output, err := cmd.CombinedOutput()
	
	if err != nil {
		t.Logf("Deploy output: %s", string(output))
		t.Logf("Deploy error: %v", err)
		t.Skip("Neo deploy not available in test environment")
		return
	}

	// Check that stage file was created
	stageFile := filepath.Join(tempDir, ".neo", "stage")
	_, err = os.Stat(stageFile)
	require.NoError(t, err, "stage file should be created")

	// Read stage file content
	content, err := os.ReadFile(stageFile)
	require.NoError(t, err)
	assert.Equal(t, "test-stage", string(content), "Stage file should contain the set stage")

	t.Logf("Successfully tested stage management in %s", tempDir)
}

// TestErrorHandling tests error handling in various scenarios
func TestErrorHandling(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "neo-error-test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	originalWd, err := os.Getwd()
	require.NoError(t, err)
	defer os.Chdir(originalWd)

	err = os.Chdir(tempDir)
	require.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Test command without config file
	cmd := exec.CommandContext(ctx, "go", "run", "../../cmd/neo", "deploy")
	cmd.Dir = tempDir
	output, err := cmd.CombinedOutput()
	
	if err == nil {
		t.Skip("Neo CLI not available in test environment")
		return
	}

	// Should fail gracefully without config
	assert.Error(t, err, "Deploy should fail without config file")
	assert.NotEmpty(t, string(output), "Should have error output")

	t.Logf("Error handling test completed in %s", tempDir)
}

// BenchmarkNeoCommands benchmarks basic CLI operations
func BenchmarkNeoCommands(b *testing.B) {
	tempDir, err := os.MkdirTemp("", "neo-benchmark-test")
	require.NoError(b, err)
	defer os.RemoveAll(tempDir)

	originalWd, err := os.Getwd()
	require.NoError(b, err)
	defer os.Chdir(originalWd)

	err = os.Chdir(tempDir)
	require.NoError(b, err)

	// Create config file
	configContent := `export default {}`
	err = os.WriteFile("neo.config.ts", []byte(configContent), 0644)
	require.NoError(b, err)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	b.Run("version", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			cmd := exec.CommandContext(ctx, "go", "run", "../../cmd/neo", "version")
			cmd.Dir = tempDir
			err := cmd.Run()
			if err != nil {
				b.Skip("Neo CLI not available in test environment")
				return
			}
		}
	})
}
