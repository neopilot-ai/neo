/// <reference path="./.neo/platform/config.d.ts" />

export default $config({
  app(input) {
    return {
      name: "aws-ai-video-generation",
      removal: input?.stage === "production" ? "retain" : "remove",
      home: "aws",
    };
  },
  async run() {
    // Note: Replace YOUR_HUGGINGFACE_TOKEN with your actual token
    const setupScript = \`
cd svd_project

# Clone the repository
git clone https://github.com/Stability-AI/generative-models.git
cd generative-models
pip install -r requirements/pt2.txt
pip install .
pip install -e git+https://github.com/Stability-AI/datapipelines.git@main#egg=sdata

huggingface-cli login --token YOUR_HUGGINGFACE_TOKEN
huggingface-cli download stabilityai/sv4d --include sv4d.safetensors --cache-dir cache
huggingface-cli download stabilityai/sv3d --include sv3d_u.safetensors --cache-dir cache

mkdir -p checkpoints
mv cache/models--stabilityai--sv4d/blobs/bdfe5bb33dfc771fc102891883befcf061873f4a96fa602037a964beca83cb44 checkpoints/sv4d.safetensors
mv cache/models--stabilityai--sv3d/blobs/d2c281b817232c492f6db27c9ce597b543187c52229cbad2a3c78e238b06c809 checkpoints/sv3d_u.safetensors

pip install --force-reinstall -v "numpy==1.25.2"
pip install imageio-ffmpeg
python scripts/sampling/simple_video_sample_4d.py --input_path assets/sv4d_videos/test_video1.mp4 --output_folder outputs/sv4d
\`;
    
    console.log("Setup script generated. Replace YOUR_HUGGINGFACE_TOKEN with your actual token.");
    return setupScript;
  }
});
