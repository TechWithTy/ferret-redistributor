import sys
from types import SimpleNamespace
from pathlib import Path
import unittest

REPO_ROOT = Path(__file__).resolve().parents[2]
if str(REPO_ROOT) not in sys.path:
    sys.path.append(str(REPO_ROOT))

from shared.tools.img import pipeline as img_pipeline  # type: ignore  # noqa: E402
from shared.tools.video import run_editly as video_tool  # type: ignore  # noqa: E402


class TestImgPipeline(unittest.TestCase):
    def test_build_command_resize_with_filter(self) -> None:
        spec = {
            "source_path": "input.jpg",
            "output_path": "output.jpg",
            "colorspace": "rgb",
            "depth": 16,
            "operations": [
                {
                    "op_name": "ResizeImage",
                    "geometry": {"width": 800, "height": 600},
                    "filter": "lanczos",
                }
            ],
        }
        cmd = img_pipeline.build_command(spec, "magick")
        self.assertEqual(
            cmd,
            [
                "magick",
                "-colorspace",
                "RGB",
                "-depth",
                "16",
                "input.jpg",
                "-filter",
                "Lanczos",
                "-resize",
                "800x600",
                "output.jpg",
            ],
        )

    def test_operation_blur(self) -> None:
        args = img_pipeline.operation_to_args({"op_name": "BlurImage", "radius": 2, "sigma": 1})
        self.assertEqual(args, ["-blur", "2x1"])

    def test_geometry_required(self) -> None:
        with self.assertRaises(SystemExit):
            img_pipeline.ensure_geometry_wh({})


class TestVideoTool(unittest.TestCase):
    def test_build_command(self) -> None:
        ns = SimpleNamespace(
            config="spec.json5",
            output="render.mp4",
            width=720,
            height=1280,
            fast=True,
        )
        cmd = video_tool.build_command("pnpm", ns)
        self.assertEqual(
            cmd,
            [
                "pnpm",
                "video:render",
                "--",
                "--config",
                "spec.json5",
                "--output",
                "render.mp4",
                "--width",
                "720",
                "--height",
                "1280",
                "--fast",
            ],
        )


if __name__ == "__main__":
    unittest.main()



